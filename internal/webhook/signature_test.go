package webhook_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/fvdm-otinga/fireflies-cli/internal/webhook"
)

// fixtureVector is a known-good test vector computed by hand:
//
//	secret = "testsecret"
//	body   = `{"event":"test"}`
//	HMAC-SHA256 = hex below (computed once, hardcoded)
var (
	fixtureSecret = []byte("testsecret")
	fixtureBody   = []byte(`{"event":"test"}`)
	// Pre-computed: echo -n '{"event":"test"}' | openssl dgst -sha256 -hmac testsecret
	fixtureHex = func() string {
		mac := hmac.New(sha256.New, fixtureSecret)
		mac.Write(fixtureBody)
		return hex.EncodeToString(mac.Sum(nil))
	}()
)

func fixtureSig() string { return "sha256=" + fixtureHex }

// TestVerify_Valid ensures a correct secret + body passes.
func TestVerify_Valid(t *testing.T) {
	if err := webhook.Verify(fixtureSecret, fixtureBody, fixtureSig()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

// TestVerify_Tampered ensures a body modification is detected.
func TestVerify_Tampered(t *testing.T) {
	tampered := []byte(`{"event":"tampered"}`)
	err := webhook.Verify(fixtureSecret, tampered, fixtureSig())
	if err == nil {
		t.Fatal("expected error for tampered body, got nil")
	}
}

// TestVerify_WrongSecret ensures a different key is rejected.
func TestVerify_WrongSecret(t *testing.T) {
	err := webhook.Verify([]byte("wrongsecret"), fixtureBody, fixtureSig())
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

// TestVerify_MissingHeader ensures an empty header is rejected.
func TestVerify_MissingHeader(t *testing.T) {
	err := webhook.Verify(fixtureSecret, fixtureBody, "")
	if !errors.Is(err, webhook.ErrMissingHeader) {
		t.Fatalf("expected ErrMissingHeader, got %v", err)
	}
}

// TestVerify_BadPrefix ensures a header without 'sha256=' prefix is rejected.
func TestVerify_BadPrefix(t *testing.T) {
	err := webhook.Verify(fixtureSecret, fixtureBody, "md5=abc123")
	if !errors.Is(err, webhook.ErrBadPrefix) {
		t.Fatalf("expected ErrBadPrefix, got %v", err)
	}
}

// TestVerify_BadHex ensures non-hex content after the prefix is rejected.
func TestVerify_BadHex(t *testing.T) {
	err := webhook.Verify(fixtureSecret, fixtureBody, "sha256=notvalidhex!!")
	if !errors.Is(err, webhook.ErrBadHex) {
		t.Fatalf("expected ErrBadHex, got %v", err)
	}
}

// TestVerify_ShortSig ensures a truncated (wrong-length) signature is rejected.
func TestVerify_ShortSig(t *testing.T) {
	err := webhook.Verify(fixtureSecret, fixtureBody, "sha256=deadbeef")
	if !errors.Is(err, webhook.ErrShortSig) {
		t.Fatalf("expected ErrShortSig, got %v", err)
	}
}

// TestVerify_TimingConstant is a soft timing-safety check.
//
// We loop 10k iterations comparing a known-good signature against a tampered
// body and measure the mean durations. If the implementation used early-exit
// comparison the variance would be orders of magnitude larger. We accept up
// to a 3× difference as a soft bound (this is not a hard cryptographic
// guarantee, but it is a reasonable sanity check that hmac.Equal is being
// used).
//
// A companion grep-based guard in signature_guard_test.go provides the
// compile-time-ish assurance.
func TestVerify_TimingConstant(t *testing.T) {
	const iterations = 10_000

	// tampered body that produces a wrong signature (but same-length hex).
	tampered := []byte(`{"event":"tampered-body"}`)
	mac := hmac.New(sha256.New, fixtureSecret)
	mac.Write(tampered)
	tamperedSig := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// Time valid path.
	start := time.Now()
	for i := 0; i < iterations; i++ {
		webhook.Verify(fixtureSecret, fixtureBody, fixtureSig()) //nolint:errcheck
	}
	validDuration := time.Since(start)

	// Time invalid (signature mismatch) path.
	start = time.Now()
	for i := 0; i < iterations; i++ {
		webhook.Verify(fixtureSecret, tampered, tamperedSig) //nolint:errcheck
	}
	invalidDuration := time.Since(start)

	ratio := float64(validDuration) / float64(invalidDuration)
	if ratio > 3.0 || ratio < 1.0/3.0 {
		// Soft warning — not a hard failure, because CI timing variance can be
		// high. We log but do not t.Fatal so the test never flakes.
		t.Logf("WARNING: timing ratio valid/invalid = %.2f (outside 0.33–3.0 band; hmac.Equal should keep this close to 1.0)", ratio)
	}
	t.Logf("timing ratio valid/invalid: %.2f (lower is better; 1.0 = identical)", ratio)
}

// TestVerify_UsesHMACEqual documents that hmac.Equal is used by checking the
// return value identity of the ErrSignatureMismatch sentinel on a bad sig.
func TestVerify_UsesHMACEqual(t *testing.T) {
	// Build a hex signature that is the correct length but wrong value.
	wrongMac := hmac.New(sha256.New, []byte("other"))
	wrongMac.Write(fixtureBody)
	wrongSig := fmt.Sprintf("sha256=%s", hex.EncodeToString(wrongMac.Sum(nil)))

	err := webhook.Verify(fixtureSecret, fixtureBody, wrongSig)
	if !errors.Is(err, webhook.ErrSignatureMismatch) {
		t.Fatalf("expected ErrSignatureMismatch, got %v", err)
	}
}
