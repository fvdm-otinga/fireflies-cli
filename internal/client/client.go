// Package client implements the Fireflies GraphQL HTTP client with auth,
// rate-limit aware retry, and per-operation token buckets.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

const DefaultEndpoint = "https://api.fireflies.ai/graphql"

// maxResponseBytes caps the size of a GraphQL response we will read into
// memory. 50 MiB is generous for real payloads but prevents a hostile or
// malfunctioning endpoint from exhausting memory.
const maxResponseBytes = 50 << 20

// maxRetryAfter bounds the delay we will honour from a Retry-After header.
// A malicious or misconfigured server could otherwise wedge the CLI for
// arbitrarily long.
const maxRetryAfter = 5 * time.Minute

// Client wraps a Fireflies GraphQL endpoint with auth + rate-limit handling.
type Client struct {
	endpoint string
	apiKey   string
	http     *http.Client
	buckets  map[string]*bucket
}

// bucket is a minimal token bucket for per-operation rate limiting.
type bucket struct {
	capacity   int
	refillTime time.Duration
	tokens     int
	next       time.Time
	mu         sync.Mutex
}

func newBucket(capacity int, refill time.Duration) *bucket {
	return &bucket{capacity: capacity, refillTime: refill, tokens: capacity}
}

func (b *bucket) take() {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	if now.After(b.next) {
		b.tokens = b.capacity
		b.next = now.Add(b.refillTime)
	}
	if b.tokens <= 0 {
		sleep := time.Until(b.next)
		if sleep > 0 {
			time.Sleep(sleep)
		}
		b.tokens = b.capacity
		b.next = time.Now().Add(b.refillTime)
	}
	b.tokens--
}

// Per-op rate limits (§2 of the plan).
var opLimits = map[string]struct {
	capacity int
	window   time.Duration
}{
	"shareMeeting":     {10, time.Hour},
	"deleteTranscript": {10, time.Minute},
	"addToLiveMeeting": {3, 20 * time.Minute},
}

// New constructs a client from a config.Profile.
func New(p config.Profile) *Client {
	return NewWithTransport(p, nil)
}

// NewWithTransport constructs a client from a config.Profile with a custom
// http.RoundTripper. If transport is nil, http.DefaultTransport is used.
// This is intended for testing only.
func NewWithTransport(p config.Profile, transport http.RoundTripper) *Client {
	endpoint := p.Endpoint
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	buckets := map[string]*bucket{}
	for op, lim := range opLimits {
		buckets[op] = newBucket(lim.capacity, lim.window)
	}
	// Default transport enforces TLS 1.2 minimum. Tests may inject their
	// own RoundTripper which overrides this.
	var defaultTransport http.RoundTripper
	if base, ok := http.DefaultTransport.(*http.Transport); ok {
		tr := base.Clone()
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		} else {
			tr.TLSClientConfig = tr.TLSClientConfig.Clone()
			tr.TLSClientConfig.MinVersion = tls.VersionTLS12
		}
		defaultTransport = tr
	} else {
		defaultTransport = http.DefaultTransport
	}
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: defaultTransport,
	}
	if transport != nil {
		httpClient.Transport = transport
	}
	return &Client{
		endpoint: endpoint,
		apiKey:   p.APIKey,
		http:     httpClient,
		buckets:  buckets,
	}
}

// Endpoint returns the resolved GraphQL endpoint.
func (c *Client) Endpoint() string { return c.endpoint }

// MakeRequest implements graphql.Client from genqlient.
// It adds the Authorization header, applies per-op rate-limiting, and
// retries on 429/503 with exponential backoff up to 3 attempts.
func (c *Client) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	if bkt, ok := c.buckets[req.OpName]; ok {
		bkt.take()
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	const maxAttempts = 3
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("build request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
		httpReq.Header.Set("User-Agent", "fireflies-cli")

		httpResp, err := c.http.Do(httpReq)
		if err != nil {
			lastErr = err
			continue
		}
		b, rerr := io.ReadAll(io.LimitReader(httpResp.Body, maxResponseBytes))
		_ = httpResp.Body.Close()
		if rerr != nil {
			lastErr = rerr
			continue
		}
		switch httpResp.StatusCode {
		case http.StatusOK:
			return json.Unmarshal(b, resp)
		case http.StatusUnauthorized, http.StatusForbidden:
			return ferr.Auth(fmt.Sprintf("unauthorized (%d): %s", httpResp.StatusCode, shortBody(b)))
		case http.StatusNotFound:
			return ferr.NotFound(string(b))
		case http.StatusTooManyRequests:
			if attempt == maxAttempts {
				return ferr.RateLimit(string(b))
			}
			wait := parseRetryAfter(httpResp.Header.Get("Retry-After"), time.Duration(attempt)*time.Second)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		case http.StatusInternalServerError, http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
			if attempt == maxAttempts {
				return ferr.General(fmt.Sprintf("server error %d: %s", httpResp.StatusCode, shortBody(b)))
			}
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		default:
			return ferr.General(fmt.Sprintf("http %d: %s", httpResp.StatusCode, shortBody(b)))
		}
	}
	if lastErr != nil {
		return ferr.General(lastErr.Error())
	}
	return ferr.General("request failed after retries")
}

func parseRetryAfter(h string, def time.Duration) time.Duration {
	d := def
	if h != "" {
		if n, err := strconv.Atoi(h); err == nil {
			d = time.Duration(n) * time.Second
		} else if t, err := http.ParseTime(h); err == nil {
			if until := time.Until(t); until > 0 {
				d = until
			}
		}
	}
	if d > maxRetryAfter {
		d = maxRetryAfter
	}
	return d
}

func shortBody(b []byte) string {
	if len(b) > 200 {
		return string(b[:200]) + "..."
	}
	return string(b)
}
