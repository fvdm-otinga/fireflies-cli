// Package webhook implements HMAC-SHA256 signature verification for
// Fireflies webhook payloads (V1 and V2).
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

// Signature verification errors — each condition has a distinct message.
var (
	ErrMissingHeader   = errors.New("webhook: x-hub-signature header is missing")
	ErrBadPrefix       = errors.New("webhook: signature header must start with 'sha256='")
	ErrBadHex          = errors.New("webhook: signature is not valid hex")
	ErrShortSig        = errors.New("webhook: signature length mismatch (truncated or empty)")
	ErrSignatureMismatch = errors.New("webhook: signature mismatch")
)

// Verify validates that signatureHeader matches HMAC-SHA256(secret, body).
//
// The header format must be "sha256=<hex>". Comparison is done via
// hmac.Equal (constant-time) to prevent timing-oracle attacks.
func Verify(secret []byte, body []byte, signatureHeader string) error {
	if signatureHeader == "" {
		return ErrMissingHeader
	}
	if !strings.HasPrefix(signatureHeader, "sha256=") {
		return ErrBadPrefix
	}
	hexSig := strings.TrimPrefix(signatureHeader, "sha256=")
	gotSig, err := hex.DecodeString(hexSig)
	if err != nil {
		return ErrBadHex
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(body)
	expectedSig := mac.Sum(nil)

	if len(gotSig) != len(expectedSig) {
		return ErrShortSig
	}

	if !hmac.Equal(gotSig, expectedSig) {
		return ErrSignatureMismatch
	}
	return nil
}
