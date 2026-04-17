# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Added

**Phase 0 — Interface Contract & Scaffold (Team Architect)**
- Repository scaffolded at `github.com/fvdm-otinga/fireflies-cli` (MIT, public)
- Interface contract (`docs/interface-contract.md`, `internal/contract/`) covering `GraphQLClient`, `OutputRenderer`, `ConfigLoader`, `ErrorModel`, and `SharedFlags`
- `genqlient` code generation configured against frozen `internal/graphql/schema.graphql`
- Smoke command `fireflies users whoami` (proves client + config + renderer integration)
- Contract acceptance tests in `internal/contract/contract_test.go`
- Exit codes: 0 success, 1 general, 2 usage, 3 auth, 4 rate-limit, 5 not-found
- `--fields` client-side projection via `internal/output/projection.go`
- Dynamic query builder stub for fat-payload commands (`internal/graphql/dynamic/`)

**Phase 1 — Team Infra**
- `internal/client/client.go`: HTTP transport, bearer-auth header, retry with exponential back-off (max 3), rate-limit token buckets (60/min global; per-op caps for `shareMeeting`, `deleteTranscript`, `addToLiveMeeting`)
- `internal/config/`: TOML profile store at `~/.config/fireflies/config.toml` (mode 0600); env `FIREFLIES_API_KEY` overrides file; `--profile` flag selects named profile
- `cmd/auth/`: `login`, `logout`, `status` subcommands
- `cmd/config/`: `get`, `set`, `list` subcommands
- `cmd/version/`: version command with `Version`, `Commit`, `Date` injected via ldflags
- `internal/output/renderer.go`: 6 output formats (`table`, `json`, `ndjson`, `yaml`, `tsv`, `plaintext`), JQ via embedded `gojq`
- `internal/output/table.go`: tablewriter-based table renderer with stable column order from `internal/output/columns/`
- `internal/pagination/`: cursor helpers for offset-based pagination
- `internal/errors/`: typed error handling, `HandleError` emits JSON on stderr

**Phase 1 — Team Read**
- `fireflies users whoami` (hardened), `users list`, `users groups` — ops: `user`, `users`, `user_groups`
- `fireflies meetings list`, `meetings get`, `meetings active` — ops: `transcripts`, `transcript`, `active_meetings`; `--transcript` flag controls sentence fetch depth
- `fireflies channels list`, `channels get` — ops: `channels`, `channel`
- `fireflies analytics` — op: `analytics`; supports `--since --until`
- `fireflies transcript text <id>` — op: `transcript.sentences`; speaker-attributed plaintext streaming
- `fireflies soundbites list`, `soundbites get` — ops: `bites`, `bite`
- `fireflies apps list` — op: `apps`
- `fireflies askfred threads`, `askfred thread <id>` — ops: `askfred_threads`, `askfred_thread`
- `fireflies rules executions <meeting-id>` — op: `rule_executions_by_meeting`
- `fireflies contacts list` — op: `contacts`
- `fireflies live items <meeting-id>` — op: `live_action_items`

**Phase 1 — Team Write**
- `fireflies meetings upload` — `uploadAudio` (URL) or `createUploadUrl` → S3 PUT → `confirmUpload` (local file); progress bar
- `fireflies meetings update title|privacy|state` — ops: `updateMeetingTitle`, `updateMeetingPrivacy`, `updateMeetingState`
- `fireflies meetings move` — op: `updateMeetingChannel`
- `fireflies users set-role` — op: `setUserRole` (destructive, requires `--yes`)
- `fireflies meetings share` — op: `shareMeeting` (rate-limited 10/hr)
- `fireflies meetings revoke` — op: `revokeSharedMeetingAccess` (destructive)
- `fireflies meetings delete` — op: `deleteTranscript` (destructive, rate-limited 10/min)
- `fireflies soundbites create` — op: `createBite`
- `fireflies askfred ask`, `askfred continue`, `askfred delete` — ops: `createAskFredThread`, `continueAskFredThread`, `deleteAskFredThread`

**Phase 1 — Team Realtime**
- `internal/realtime/client.go`: Socket.IO over WSS to `wss://api.fireflies.ai/ws/realtime`, auth handshake, reconnect
- `fireflies realtime tail <meeting-id>` — streams transcript events; supports `--output json|plaintext`
- `fireflies live add`, `live soundbite`, `live action-item` — ops: `addToLiveMeeting`, `createLiveSoundbite`, `createLiveActionItem`
- `fireflies webhooks serve` — HTTP receiver for V1 + V2 webhooks; HMAC-SHA256 via `hmac.Equal` (constant-time); emits verified events as NDJSON on stdout

**Phase 1 — Team Release (this work)**
- `.goreleaser.yaml`: 5 targets (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64); tar.gz / zip archives; `checksums.txt`; CycloneDX SBOM via `syft` (optional); changelog from conventional commits; Homebrew tap formula auto-published to `fvdm-otinga/homebrew-tap`
- `.github/workflows/release.yml`: triggers on `v*` tags, runs goreleaser
- `.github/workflows/ci.yml`: adds `completion-zsh` (syntax check) and `fixture-secrets-scan` (grep for leaked credentials) jobs
- `cmd/completion.go`: explicit completion command (bash/zsh/fish/powershell) with installation instructions in `--help`
- `docs/install.md`: install instructions for Homebrew, `go install`, and manual binary download
- `docs/command-reference.md`: annotated command tree and exit codes
- `docs/op-command-map.md`: all 35 GraphQL operations mapped to CLI commands
- `docs/reference/`: auto-generated per-command Markdown via `make docs`
- `scripts/gen-docs.go` (`//go:build gendocs`): Cobra doc generator
- `scripts/scrub-fixtures.sh`: fixture scrubber/checker for CI and dev use
- `Makefile`: `build`, `test`, `lint`, `generate`, `docs`, `release-dry`, `scrub` targets
- `README.md`: full documentation update
- `CHANGELOG.md`: this file (Keep-a-Changelog format)
- `fvdm-otinga/homebrew-tap` GitHub repo bootstrapped

---

[Unreleased]: https://github.com/fvdm-otinga/fireflies-cli/compare/contract-v1...HEAD
