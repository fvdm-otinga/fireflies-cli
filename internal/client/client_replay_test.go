// Package client_test — go-vcr cassette recording and replay tests.
//
// Recording: set TEST_RECORD=1 and FIREFLIES_API_KEY to record a new cassette.
// Replay (default, CI-safe): runs against the committed cassette with no network.
//
// Example — record:
//
//	TEST_RECORD=1 FIREFLIES_API_KEY=<key> go test ./internal/client/ -run TestWhoami_Replay -v
//
// The cassette is saved to testdata/fixtures/whoami_happy.yaml and must be
// scrubbed before committing:
//
//	scripts/scrub-fixtures.sh testdata/fixtures/
package client_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// fixturePath returns the absolute path to a fixture file, anchored to the
// repo root regardless of the test binary's working directory.
func fixturePath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	// thisFile = .../internal/client/client_replay_test.go
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	return filepath.Join(repoRoot, "testdata", "fixtures", name)
}

// whoamiGraphQLBody is the GraphQL request body for the Whoami query.
// The query uses literal newlines to match what the client would send.
const whoamiGraphQLBody = "{\"query\":\"query Whoami {\\n  user {\\n    user_id\\n    email\\n    name\\n    is_admin\\n    num_transcripts\\n    minutes_consumed\\n  }\\n}\",\"operationName\":\"Whoami\"}"

const firefliesEndpoint = "https://api.fireflies.ai/graphql"

// TestWhoami_Replay replays the whoami_happy.yaml cassette without hitting the
// network. If TEST_RECORD=1, it makes a live request and records a new cassette.
func TestWhoami_Replay(t *testing.T) {
	// go-vcr appends ".yaml" automatically; pass the path without extension.
	fixtureBase := fixturePath("whoami_happy")
	fixtureYAML := fixtureBase + ".yaml"

	mode := recorder.ModeReplayOnly
	if os.Getenv("TEST_RECORD") == "1" {
		mode = recorder.ModeRecordOnce
	}

	// In replay mode, skip the test if the cassette does not exist yet.
	if mode == recorder.ModeReplayOnly {
		if _, err := os.Stat(fixtureYAML); os.IsNotExist(err) {
			t.Skip("fixture not recorded yet — run with TEST_RECORD=1 to record")
		}
	}

	// Custom matcher: match on method + URL only.
	// Body and auth headers are ignored so the cassette is replay-safe with any key
	// and regardless of minor request-body whitespace differences.
	urlMatcher := cassette.MatcherFunc(func(r *http.Request, i cassette.Request) bool {
		return r.Method == i.Method && r.URL.String() == i.URL
	})

	r, err := recorder.New(
		fixtureBase,
		recorder.WithMode(mode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithMatcher(urlMatcher),
	)
	if err != nil {
		t.Fatalf("recorder.New: %v", err)
	}
	defer func() {
		if err := r.Stop(); err != nil {
			t.Errorf("recorder.Stop: %v", err)
		}
	}()

	apiKey := os.Getenv("FIREFLIES_API_KEY")
	if apiKey == "" {
		apiKey = "REDACTED" // replay mode uses cassette; key value doesn't matter
	}

	// Use the VCR-instrumented http.Client directly.
	httpClient := r.GetDefaultClient()

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		firefliesEndpoint,
		bytes.NewBufferString(whoamiGraphQLBody),
	)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", "fireflies-cli-test")

	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	// Parse and validate the response structure.
	var gqlResp struct {
		Data struct {
			User *struct {
				UserID *string `json:"user_id"`
				Email  *string `json:"email"`
				Name   *string `json:"name"`
			} `json:"user"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(bodyBytes, &gqlResp); err != nil {
		t.Fatalf("decode response: %v\nbody: %s", err, string(bodyBytes))
	}

	if len(gqlResp.Errors) > 0 {
		msgs := make([]string, len(gqlResp.Errors))
		for i, e := range gqlResp.Errors {
			msgs[i] = e.Message
		}
		t.Fatalf("GraphQL errors: %s", strings.Join(msgs, "; "))
	}

	if gqlResp.Data.User == nil {
		t.Fatal("expected data.user to be non-nil")
	}
	if gqlResp.Data.User.UserID == nil || *gqlResp.Data.User.UserID == "" {
		t.Error("expected non-empty user_id")
	}

	t.Logf("whoami: user_id=%v email=%v name=%v",
		ptrStr(gqlResp.Data.User.UserID),
		ptrStr(gqlResp.Data.User.Email),
		ptrStr(gqlResp.Data.User.Name),
	)
}

func ptrStr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
