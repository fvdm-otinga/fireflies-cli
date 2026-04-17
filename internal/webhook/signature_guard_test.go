// Package webhook_test contains a source-guard test that ensures no
// non-constant-time comparison is ever introduced into signature.go.
package webhook_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestSignatureGuard_NoDirectComparison reads signature.go and fails if it
// contains any of the forbidden non-constant-time comparison patterns.
// This is the "compile-time-ish" guard required by §9 of the plan.
func TestSignatureGuard_NoDirectComparison(t *testing.T) {
	// Locate signature.go relative to this test file.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	srcFile := filepath.Join(dir, "signature.go")

	data, err := os.ReadFile(srcFile)
	if err != nil {
		t.Fatalf("cannot read signature.go: %v", err)
	}
	src := string(data)

	forbidden := []string{
		"== sig",
		"!= sig",
		"bytes.Equal(sig",
		"bytes.Equal(expect",
	}
	for _, pattern := range forbidden {
		if strings.Contains(src, pattern) {
			t.Errorf("signature.go contains forbidden non-constant-time pattern: %q", pattern)
		}
	}

	// Positive check: hmac.Equal must be present.
	if !strings.Contains(src, "hmac.Equal(") {
		t.Error("signature.go does not call hmac.Equal — constant-time comparison is required")
	}
}
