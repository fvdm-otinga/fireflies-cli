# Security Audit — v1.0.0

Performed: 2026-04-17  
Auditor: Phase 2 integration agent

---

## Gate Results

| # | Gate | Result | Notes |
|---|---|---|---|
| 1 | `grep -r '== sig' internal/webhook/ cmd/webhooks/` finds nothing outside guard test | PASS | Only occurrence is in `internal/webhook/signature_guard_test.go` (the guard test itself) |
| 2 | `grep -r 'bytes.Equal(sig' internal/webhook/ cmd/webhooks/` finds nothing outside guard test | PASS | Only occurrence is in `internal/webhook/signature_guard_test.go` |
| 3 | `grep -rn 'Bearer ' internal/ cmd/` — no matches log the token | PASS | Bearer used only in `Authorization` header set (`internal/client/client.go:115`, `internal/realtime/client.go:73`); never written to any logger or output stream |
| 4 | `internal/config/config.go` writes config with mode `0600` | PASS | `os.WriteFile(p, b, 0600)` at line 131; dir created with `0700` |
| 5 | `FIREFLIES_API_KEY=secret-abcdef123456 ./fireflies --help` contains no literal key | PASS | Verified for root `--help` and subgroups: `users`, `meetings`, `auth`, `config`, `version` |
| 6 | `TestAuthTokenNotLogged` unit test in `internal/client/client_test.go` | PASS | Test verifies Authorization header is set correctly and secret does not appear in URL/User-Agent |
| 7 | `scripts/scrub-fixtures.sh --check testdata/fixtures/` exits 0 | PASS | No sensitive patterns found in committed fixtures |

---

## Detailed Findings

### Webhook Signature Verification
- File: `internal/webhook/signature.go`
- Uses `hmac.Equal(gotSig, expectedSig)` (constant-time comparison) — correct
- Rejects missing header, bad hex, length mismatch, and value mismatch as distinct errors
- `TestWebhookSignatureHMAC` and `TestSignatureGuard_NaiveComparisonNotPresent` both pass

### Bearer Token Handling
- `internal/client/client.go`: sets `Authorization: Bearer <key>` in outbound HTTP headers only
- `internal/realtime/client.go`: same pattern for WebSocket dial headers
- No `log.Print*` or `fmt.Fprintf(stderr)` calls include the API key string
- The `--verbose` flag does not exist (no verbose logging at all in v1.0.0)

### Config File Permissions
- `internal/config/config.go` line 131: `os.WriteFile(p, b, 0600)`
- Config directory created with mode `0700`
- Covered by `TestConfigFilePermissions` in `internal/contract/contract_test.go`

### Fixture Secrets
- `testdata/fixtures/whoami_happy.yaml`: Authorization header replaced with `[SCRUBBED]`; email field uses `[scrubbed]` (no `@` symbol); all UUIDs use placeholder values
- Scrub script patterns checked: Bearer tokens, email addresses, UUIDs, `FIREFLIES_API_KEY=`

---

## Known Limitation

`govulncheck ./...` reports 4 vulnerabilities in `crypto/x509` and `crypto/tls` at go1.26.1, all fixed in go1.26.2. The go1.26.2 toolchain is not yet available. These affect the TLS stack used by the HTTP client and WebSocket client (not application logic). **Ship blocker for v1.0.0 if go1.26.2 is released before tagging.** Monitor https://pkg.go.dev/vuln/GO-2026-4870, GO-2026-4866, GO-2026-4946, GO-2026-4947.

---

## Conclusion

All 7 required security gates pass. The single known open issue (govulncheck stdlib CVEs) is a toolchain upgrade dependency, not a code-level defect.
