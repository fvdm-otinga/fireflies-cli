// Package webhook_test contains integration tests for the webhook HTTP handler.
//
// We import the mux builder from cmd/webhooks via a shared helper here so we
// don't need to export it from that package. Instead this test reimplements
// the small mux inline to keep the webhook package self-contained.
package webhook_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fvdm-otinga/fireflies-cli/internal/webhook"
)

// minimalServeHandler mimics the cmd/webhooks handler logic so we can test
// the webhook.Verify integration without importing cmd/.
func minimalServeHandler(secret []byte, version string, out io.Writer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close() //nolint:errcheck

		sigHeader := r.Header.Get("x-hub-signature")
		if err := webhook.Verify(secret, body, sigHeader); err != nil {
			http.Error(w, "unauthorized: signature verification failed", http.StatusUnauthorized)
			return
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			payload = map[string]any{"raw": string(body)}
		}
		payload["version"] = version

		line, _ := json.Marshal(payload)
		_, _ = fmt.Fprintf(out, "%s\n", line) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}
}

func signBody(secret, body []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

// TestServe_ValidSignature verifies that a correctly signed POST returns 200
// and emits a NDJSON line on stdout.
func TestServe_ValidSignature(t *testing.T) {
	secret := []byte("testsecret")
	body := []byte(`{"event":"test","payload":"hello"}`)

	var stdout bytes.Buffer
	handler := minimalServeHandler(secret, "v2", &stdout)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/v2", bytes.NewReader(body))
	req.Header.Set("x-hub-signature", signBody(secret, body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	line := strings.TrimSpace(stdout.String())
	if line == "" {
		t.Fatal("expected NDJSON line on stdout, got empty")
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(line), &out); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if out["version"] != "v2" {
		t.Errorf("version: got %v, want v2", out["version"])
	}
	if out["event"] != "test" {
		t.Errorf("event: got %v, want test", out["event"])
	}
}

// TestServe_BadSignature verifies that a tampered signature returns 401 and
// does not emit any output.
func TestServe_BadSignature(t *testing.T) {
	secret := []byte("testsecret")
	body := []byte(`{"event":"test"}`)

	var stdout bytes.Buffer
	handler := minimalServeHandler(secret, "v2", &stdout)

	req := httptest.NewRequest(http.MethodPost, "/webhooks/v2", bytes.NewReader(body))
	req.Header.Set("x-hub-signature", "sha256=000000")
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if stdout.Len() != 0 {
		t.Errorf("expected no stdout output for bad sig, got: %s", stdout.String())
	}
}
