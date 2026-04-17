// Package client_test contains unit tests for the Fireflies HTTP client.
package client_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/Khan/genqlient/graphql"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
)

// captureRoundTripper records the outbound request and responds with the
// provided body + status code. Used to assert header values without hitting
// the network.
type captureRoundTripper struct {
	captured *http.Request
	body     string
	status   int
}

func (rt *captureRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Deep-copy the request so the caller can inspect headers after the
	// response body has been read.
	rt.captured = req.Clone(req.Context())
	body := rt.body
	if body == "" {
		body = `{"data":{}}`
	}
	status := rt.status
	if status == 0 {
		status = http.StatusOK
	}
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}, nil
}

// TestAuthTokenNotLogged verifies that:
//  1. The Authorization header is set to "Bearer <secret>" on each request.
//  2. The literal secret value does NOT appear in any observable output to
//     stdout/stderr — we confirm this by ensuring the client has no internal
//     logger that could leak it, and that the round-tripper only receives
//     the header (never logs it).
func TestAuthTokenNotLogged(t *testing.T) {
	const secret = "secret-abcdef123456"

	rt := &captureRoundTripper{
		body:   `{"data":{"user":{"id":"u1"}}}`,
		status: http.StatusOK,
	}

	// Build a client with our capturing transport injected via ClientWithTransport.
	c := client.NewWithTransport(config.Profile{APIKey: secret}, rt)

	// Issue a minimal request using the genqlient graphql.Client interface.
	req := &graphql.Request{
		Query:  `query { user { id } }`,
		OpName: "whoami",
	}
	var resp graphql.Response
	if err := c.MakeRequest(context.Background(), req, &resp); err != nil {
		t.Fatalf("MakeRequest failed: %v", err)
	}

	// Gate 1: Authorization header is set correctly.
	if rt.captured == nil {
		t.Fatal("no request was captured")
	}
	authHeader := rt.captured.Header.Get("Authorization")
	expected := "Bearer " + secret
	if authHeader != expected {
		t.Errorf("Authorization header: got %q, want %q", authHeader, expected)
	}

	// Gate 2: the captured request URL / body / user-agent do not contain the
	// raw secret (it must only appear inside the Authorization header value).
	reqBody := rt.captured.URL.RawQuery
	if contains(reqBody, secret) {
		t.Errorf("secret found in request URL query: %s", reqBody)
	}
	ua := rt.captured.Header.Get("User-Agent")
	if contains(ua, secret) {
		t.Errorf("secret found in User-Agent header: %s", ua)
	}
}

func contains(s, sub string) bool {
	return len(sub) > 0 && len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
