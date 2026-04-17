package errors_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

func TestConstructors(t *testing.T) {
	tests := []struct {
		fn       func(string) *ferr.CLIError
		msg      string
		wantCode string
		wantExit ferr.ExitCode
	}{
		{ferr.Auth, "bad key", "auth", ferr.ExitAuthError},
		{ferr.Usage, "bad flag", "usage", ferr.ExitUsageError},
		{ferr.NotFound, "missing", "not_found", ferr.ExitNotFound},
		{ferr.RateLimit, "slow down", "rate_limit", ferr.ExitRateLimit},
		{ferr.General, "oops", "error", ferr.ExitGeneralError},
	}
	for _, tc := range tests {
		e := tc.fn(tc.msg)
		if e.Error() != tc.msg {
			t.Errorf("%s: Error() = %q, want %q", tc.wantCode, e.Error(), tc.msg)
		}
		if e.Exit != tc.wantExit {
			t.Errorf("%s: Exit = %d, want %d", tc.wantCode, e.Exit, tc.wantExit)
		}
		if e.Code != tc.wantCode {
			t.Errorf("Code = %q, want %q", e.Code, tc.wantCode)
		}
	}
}

func TestNewf(t *testing.T) {
	e := ferr.Newf(ferr.KindAuth, "key %s invalid", "abc")
	if e.Exit != ferr.ExitAuthError {
		t.Errorf("exit = %d, want %d", e.Exit, ferr.ExitAuthError)
	}
	if e.Message != "key abc invalid" {
		t.Errorf("message = %q", e.Message)
	}

	e2 := ferr.Newf(ferr.KindRateLimit, "too fast: %d req/s", 100)
	if e2.Exit != ferr.ExitRateLimit {
		t.Errorf("exit = %d, want %d", e2.Exit, ferr.ExitRateLimit)
	}

	e3 := ferr.Newf(ferr.KindNotFound, "item %d", 42)
	if e3.Exit != ferr.ExitNotFound {
		t.Errorf("exit = %d, want %d", e3.Exit, ferr.ExitNotFound)
	}

	e4 := ferr.Newf(ferr.KindUsage, "bad args")
	if e4.Exit != ferr.ExitUsageError {
		t.Errorf("exit = %d, want %d", e4.Exit, ferr.ExitUsageError)
	}

	e5 := ferr.Newf(ferr.KindGeneral, "something broke")
	if e5.Exit != ferr.ExitGeneralError {
		t.Errorf("exit = %d, want %d", e5.Exit, ferr.ExitGeneralError)
	}
}

func TestFromGraphQLError(t *testing.T) {
	tests := []struct {
		input    string
		wantExit ferr.ExitCode
	}{
		{"Unauthorized: invalid token", ferr.ExitAuthError},
		{"unauthorized access", ferr.ExitAuthError},
		{"Forbidden: no permissions", ferr.ExitAuthError},
		{"rate limit exceeded", ferr.ExitRateLimit},
		{"Too many requests", ferr.ExitRateLimit},
		{"resource not found", ferr.ExitNotFound},
		{"Not Found", ferr.ExitNotFound},
		{"internal server error", ferr.ExitGeneralError},
		{"something went wrong", ferr.ExitGeneralError},
	}
	for _, tc := range tests {
		e := ferr.FromGraphQLError(fmt.Errorf("%s", tc.input))
		if e == nil {
			t.Errorf("%q: got nil", tc.input)
			continue
		}
		if e.Exit != tc.wantExit {
			t.Errorf("%q: exit = %d, want %d", tc.input, e.Exit, tc.wantExit)
		}
	}

	// nil input returns nil
	if ferr.FromGraphQLError(nil) != nil {
		t.Error("expected nil for nil input")
	}

	// CLIError passthrough
	original := ferr.Auth("already typed")
	result := ferr.FromGraphQLError(original)
	if result.Exit != ferr.ExitAuthError {
		t.Errorf("passthrough: exit = %d, want %d", result.Exit, ferr.ExitAuthError)
	}
}

func TestHandle_CLIError(t *testing.T) {
	var buf bytes.Buffer
	e := ferr.Auth("bad key")
	code := ferr.Handle(&buf, e)
	if code != ferr.ExitAuthError {
		t.Errorf("exit code = %d, want %d", code, ferr.ExitAuthError)
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Errorf("stderr not valid JSON: %v", err)
	}
	if out["code"] != "auth" {
		t.Errorf("code = %v", out["code"])
	}
}

func TestHandle_PlainError(t *testing.T) {
	var buf bytes.Buffer
	code := ferr.Handle(&buf, fmt.Errorf("plain error"))
	if code != ferr.ExitGeneralError {
		t.Errorf("exit code = %d, want %d", code, ferr.ExitGeneralError)
	}
	if buf.String() != "plain error\n" {
		t.Errorf("stderr = %q", buf.String())
	}
}

func TestHandle_Nil(t *testing.T) {
	var buf bytes.Buffer
	code := ferr.Handle(&buf, nil)
	if code != ferr.ExitSuccess {
		t.Errorf("exit code = %d, want 0", code)
	}
}
