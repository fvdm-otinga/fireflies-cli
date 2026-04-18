# Plan — `fireflies` CLI

**Goal:** Ship a Go CLI at `github.com/fvdm-otinga/fireflies-cli` that wraps the full Fireflies.ai GraphQL API (35 operations + webhooks + Socket.IO realtime) in a token-efficient, LLM-friendly way.

**Hosting note:** This project lives on **GitHub** by explicit user direction (2026-04-17), overriding the global GitLab-first default for coding projects. Architectural sibling for reference: `~/Documents/code/notion-cli` (4ier/notion-cli) — a Go CLI with multi-profile auth built for the same "CLI-over-MCP for token efficiency" goal.

**Approach:** Approach C — five parallel teams build simultaneously against a frozen Phase 0 interface contract. Phase 2 integrates and hardens; Phase 3 tags v1.0.0 and publishes.

---

## 1. Locked decisions

| Area | Decision |
|---|---|
| Language/stack | Go + Cobra + `genqlient` (type-safe codegen) + `itchyny/gojq` (embedded jq) |
| Repo | `github.com/fvdm-otinga/fireflies-cli` (public, MIT) |
| Scope | All 17 queries + 18 mutations + webhook receiver (V1+V2) + Socket.IO realtime |
| Auth | `FIREFLIES_API_KEY` env > `~/.config/fireflies/config.toml` > `fireflies auth login`; multi-profile via `--profile` |
| Default output | Human-readable tables (gh-style); `--json` opt-in |
| Extra output formats | `--output json\|ndjson\|yaml\|tsv\|plaintext`, `--jq <expr>` post-filter, `--fields a,b,c` → GraphQL selection set |
| Transcript token control | `--transcript=none\|preview\|full` (list=none, detail=preview); dedicated `fireflies transcript text <id>` command; `--since/--until` windowing; `--format plaintext` speaker-attributed |
| Distribution | GitHub Actions CI; `goreleaser` publishes darwin arm64/amd64 + linux amd64/arm64 + windows amd64 to GitHub Releases; Homebrew tap at `github.com/fvdm-otinga/homebrew-tap`; `go install` works always |
| Tests | Unit per command + `go-vcr` HTTP fixtures + `golangci-lint` + `go vet` + CI build matrix |
| Git | HTTPS only, no Claude co-author trailers |

---

## 2. Full API inventory (v1.0.0 must cover all 35 ops)

**Queries (17):** `transcript`, `transcripts`, `user`, `users`, `user_groups`, `channels`, `channel`, `active_meetings`, `bite`, `bites`, `apps`, `askfred_threads`, `askfred_thread`, `contacts`, `analytics`, `live_action_items`, `rule_executions_by_meeting`.

**Mutations (18):** `uploadAudio`, `createUploadUrl`, `confirmUpload`, `addToLiveMeeting`, `updateMeetingState`, `updateMeetingTitle`, `updateMeetingPrivacy`, `updateMeetingChannel`, `shareMeeting`, `revokeSharedMeetingAccess`, `deleteTranscript`, `createBite`, `createLiveSoundbite`, `createLiveActionItem`, `createAskFredThread`, `continueAskFredThread`, `deleteAskFredThread`, `setUserRole`.

`uploadAudio` (single-step upload of a public URL) and the `createUploadUrl` → `confirmUpload` pair (two-step upload of a local file) are three separate named mutations, all exposed by a single `fireflies meetings upload` command that picks the right flow based on input (file path vs. URL).

**Other surfaces:** Webhooks V1 + V2 (HMAC-SHA256 `x-hub-signature`); Socket.IO WebSocket at `wss://api.fireflies.ai/ws/realtime`.

**Rate limits (enforced client-side with token buckets):** 60 req/min Business tier; `shareMeeting` 10/hr; `deleteTranscript` 10/min; `addToLiveMeeting` 3 per 20 min. Pagination is offset-based (`skip`+`limit`); `transcripts.limit` caps at 50.

---

## 3. Interface contract (frozen at Phase 0 exit, tag `contract-v1`)

Document at `docs/interface-contract.md`. Package at `internal/contract/`.

```
GraphQLClient:
  Do(ctx, opName string, vars any, out any) error
  WithProfile(name string) GraphQLClient
  // handles auth header, rate-limit buckets, 429 exponential backoff (max 3 retries),
  // translates GraphQL errors to typed ErrorModel values

OutputRenderer:
  Table(headers []string, rows [][]string) error
  JSON(v any) error
  NDJSON(v any) error
  YAML(v any) error
  TSV(headers []string, rows [][]string) error
  Plaintext(lines []string) error
  JQ(v any, expr string) error   // via itchyny/gojq, no subprocess

ConfigLoader:
  Profile(name string) (*Profile, error)
  APIKey() string      // env wins over file
  Save(p Profile) error // file mode 0600

ErrorModel (exit codes):
  0 success
  1 general / API error      (stderr JSON: {"error":"...","code":"..."})
  2 usage error
  3 auth error
  4 rate limit exhausted
  5 not found

SharedFlags (bound on every command via internal/flags/flags.go):
  --profile string, --json, --jq string, --output string,
  --fields string, --limit int, --skip int,
  --transcript none|preview|full, --since RFC3339, --until RFC3339,
  --yes (confirm-bypass), --dry-run
```

The contract is import-only after Phase 0; changes require a documented RFC in `docs/rfc/` and owner approval.

### 3.1 `--fields` selection-set strategy

`genqlient` produces fixed typed query structs — it cannot narrow a query's selection set at runtime. The contract resolves this explicitly:

- **Default path (all commands):** `--fields` performs **client-side projection** after decode. The full response is fetched, then the renderer keeps only the named top-level (and dotted-nested) paths before table/JSON/JQ rendering. This saves output tokens but not API bytes. Implemented once in `internal/output/projection.go`; every command gets it for free.
- **Fat-payload path (3 commands):** `meetings list`, `meetings get`, and `transcript text` additionally ship with **hand-written dynamic query builders** under `internal/graphql/dynamic/` that assemble the outbound GraphQL document string from a selection whitelist. These commands accept `--fields` and rewrite the *outbound* query to skip `sentences`, `summary.*`, `analytics`, etc. This is the only place the CLI bypasses `genqlient` typing — typed response-side structs are still used via a shared `Transcript` union struct with pointer fields (nil = not requested).
- **Discovery:** `--fields ?` prints the whitelist for that command and exits.

### 3.2 Renderer input shape (`JQ`, `JSON`, `NDJSON`, `YAML`)

Commands always pass the **top-level `genqlient` response struct** (or a slice of them for list ops) to the renderer. The renderer marshals to JSON with `json.Marshal` using struct tags as produced by `genqlient`, then `gojq` operates on the parsed JSON value. Renderer never receives a map or any hand-built shape — this forbids per-command shape variance. A shared `internal/output/envelope.go` wraps list responses as `{"data":[...], "meta":{"limit":N, "skip":N, "next_skip":N|null}}`; single-object responses are emitted unwrapped.

### 3.3 Table column definitions

Table column order and header text for each response type live in `internal/output/columns/` (one Go file per response type: `transcript.go`, `user.go`, `channel.go`, ...) with a `Columns() []ColumnDef` function. Every command for that type imports the same column definition — there is no per-command table formatting.

---

## 4. Phase 0 — Interface Contract & Scaffold Bootstrap

**Owner:** 1 architect agent (`general-purpose`, sonnet). Sequential, blocks Phase 1.

### Tasks

1. Create repo `github.com/fvdm-otinga/fireflies-cli` (public, MIT), clone over HTTPS
2. `go mod init github.com/fvdm-otinga/fireflies-cli`
3. Scaffold directories: `cmd/`, `internal/{client,config,output,pagination,errors,flags,graphql,contract,realtime,webhook}/`, `docs/`, `testdata/fixtures/`, `benchmarks/token_efficiency/`
4. Fetch Fireflies GraphQL schema via introspection (live API key required); fallback to hand-derived `schema.graphql` if introspection is blocked. Commit to `internal/graphql/schema.graphql`.
5. Configure `genqlient.yaml`; verify `go generate ./...` produces typed bindings from a placeholder query.
6. Write interface contract (`docs/interface-contract.md` + `internal/contract/*.go`) covering everything in §3 including the `--fields` strategy (§3.1), renderer input shape (§3.2), and column definitions (§3.3). Create `internal/output/projection.go` with the client-side field projector and one `internal/graphql/dynamic/transcript.go` stub proving the dynamic builder pattern works. Create contract acceptance tests in `internal/contract/contract_test.go` (see §8).
7. Implement ONE smoke command end-to-end: `fireflies users whoami` — proves client, config, output renderer, and error model integrate cleanly.
8. Repo hygiene: `.gitignore`, `.github/CODEOWNERS`, `README.md` stub, `LICENSE` (MIT), `Makefile` (targets: `build test lint generate release-dry`).
9. Stub `.github/workflows/ci.yml` (jobs defined, bodies empty — Team Release fills in Phase 1).
10. Tag `contract-v1` on the Phase 0 merge commit.

### Exit gates

- [ ] `internal/contract/contract_test.go` passes (all GraphQLClient, OutputRenderer, ConfigLoader, ErrorModel, FlagSet tests)
- [ ] `fireflies users whoami` returns valid JSON with `--json` against a live API key
- [ ] `go generate ./...` is deterministic (CI-safe)
- [ ] Contract tag `contract-v1` pushed

**Effort:** 3–4 agent-hours.

---

## 5. Phase 1 — Parallel Build (5 teams)

**Entry:** `contract-v1` tag exists; all teams branch from there.

Each team lands its branch via PR against `main`. No team modifies `internal/contract/` or `internal/graphql/schema.graphql`.

### Team Infra (2 agents) — branch `infra/core`

| Agent | Scope |
|---|---|
| A1 | `internal/client/client.go` (HTTP, auth, retry, rate-limit token buckets), `internal/config/` (TOML load/save, env override, profiles), `cmd/auth/{login,logout,status}.go`, `cmd/config/{get,set,list}.go`, `cmd/version.go` |
| A2 | `internal/output/renderer.go` (6 formats + JQ via gojq), `internal/output/table.go` (`tablewriter`-based), `internal/pagination/`, `internal/errors/` (typed errors + `HandleError(err)` emits stderr JSON), unit tests for renderer + pagination |

### Team Read (2 agents) — branch `read/queries`

Implements 17 read commands. One `.graphql` file per op in `internal/graphql/queries/`.

**Agent R1:**

| Command | GraphQL op |
|---|---|
| `fireflies users whoami` (harden Phase-0 stub) | `user` self |
| `fireflies users list` | `users` |
| `fireflies users groups` | `user_groups` |
| `fireflies meetings list` | `transcripts` (supports `--since --until --limit --skip --transcript`) |
| `fireflies meetings get <id>` | `transcript` (supports `--transcript --fields`) |
| `fireflies meetings active` | `active_meetings` |
| `fireflies channels list` | `channels` |
| `fireflies channels get <id>` | `channel` |
| `fireflies analytics` | `analytics` (`--since --until`) |

**Agent R2:**

| Command | GraphQL op |
|---|---|
| `fireflies transcript text <id>` | `transcript.sentences` streamed as `Speaker: text` plaintext; supports `--since --until` |
| `fireflies soundbites list` | `bites` |
| `fireflies soundbites get <id>` | `bite` |
| `fireflies apps list` | `apps` |
| `fireflies askfred threads` | `askfred_threads` |
| `fireflies askfred thread <id>` | `askfred_thread` |
| `fireflies rules executions <meeting-id>` | `rule_executions_by_meeting` |
| `fireflies contacts list` | `contacts` |
| `fireflies live items <meeting-id>` | `live_action_items` |

### Team Write (2 agents) — branch `write/mutations`

All destructive ops require `--yes` or an interactive confirm. All mutations support `--dry-run` (prints GraphQL op body + variables, zero HTTP).

**Agent W1:**

| Command | Mutation |
|---|---|
| `fireflies meetings upload <file-or-url>` | `createUploadUrl` → S3 PUT → `confirmUpload` (local file path), or `uploadAudio` (public URL); progress bar |
| `fireflies meetings update title <id> <title>` | `updateMeetingTitle` |
| `fireflies meetings update privacy <id> <level>` | `updateMeetingPrivacy` |
| `fireflies meetings update state <id> <state>` | `updateMeetingState` |
| `fireflies meetings move <id> --channel <cid>` | `updateMeetingChannel` |
| `fireflies users set-role <uid> <role>` | `setUserRole` (destructive → `--yes`) |

**Agent W2:**

| Command | Mutation |
|---|---|
| `fireflies meetings share <id> --email <e>...` | `shareMeeting` (rate bucket 10/hr) |
| `fireflies meetings revoke <id> --email <e>` | `revokeSharedMeetingAccess` (destructive) |
| `fireflies meetings delete <id>` | `deleteTranscript` (destructive, rate 10/min) |
| `fireflies soundbites create --meeting <id> --start t --end t` | `createBite` |
| `fireflies askfred ask --meeting <id> <question>` | `createAskFredThread` |
| `fireflies askfred continue <thread-id> <question>` | `continueAskFredThread` |
| `fireflies askfred delete <thread-id>` | `deleteAskFredThread` (destructive) |

### Team Realtime (1 agent) — branch `realtime/ws`

1. `internal/realtime/client.go` — Socket.IO over WSS to `wss://api.fireflies.ai/ws/realtime`, auth handshake, reconnect. Library candidate: `github.com/graarh/golang-socketio` (fallback: raw `nhooyr.io/websocket` + hand-rolled Socket.IO frames).
2. `cmd/realtime/tail.go` — `fireflies realtime tail <meeting-id>` streams transcript events; supports `--output json|plaintext`.
3. `cmd/live/add.go` — `addToLiveMeeting` (rate 3/20min).
4. `cmd/live/soundbite.go` — `createLiveSoundbite`.
5. `cmd/live/action-item.go` — `createLiveActionItem`.
6. `cmd/webhooks/serve.go` — `fireflies webhooks serve --port 8080` HTTP receiver (secret via `FIREFLIES_WEBHOOK_SECRET` env or `--secret-stdin`), V1 + V2 routes, HMAC-SHA256 via `hmac.Equal` (constant-time), emits verified events as NDJSON on stdout. Hardened with 1 MiB body cap, server timeouts, and panic-recovery middleware.

### Team Release (1 agent) — branch `release/infra`

1. `.goreleaser.yaml` — 5 targets, archives, `checksums.txt`, CycloneDX SBOM, changelog from tags.
2. `.github/workflows/ci.yml` — jobs `lint` (`golangci-lint v1.57`), `vet` (`go vet ./...`), `test` (`go test -race ./...` across `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`), `build` (matrix). Windows builds but is not in test matrix.
3. `.github/workflows/release.yml` — triggers on `v*` tag, runs goreleaser.
4. `.golangci.yml` — `errcheck`, `gosimple`, `staticcheck`, `unused`, `govet`, `ineffassign`.
5. `Formula/fireflies.rb` stub in the separate `homebrew-tap` repo.
6. Shell completions via Cobra built-in (`fireflies completion bash|zsh|fish|powershell`).
7. `docs/install.md`.

### Phase 1 exit gates

- [ ] All 5 branches build green in CI against `contract-v1`
- [ ] Every query command returns valid JSON under `--json`
- [ ] Every mutation command implements `--dry-run`; destructive ones honour `--yes`
- [ ] `--fields a,b,c` narrows the outbound GraphQL selection set (verified via go-vcr request-body inspection)
- [ ] Rate-limit token buckets enforce documented per-op caps (unit-tested with fake clock)
- [ ] Coverage ≥ 70 % in `internal/client`, `internal/config`, `internal/output`
- [ ] Webhook signature verification uses `hmac.Equal`; `TestWebhookSignatureTimingConstant` present
- [ ] No API keys or real PII in any `testdata/fixtures/*.yaml` (scrub check passes)

**Effort:** 12–16 agent-hours total, ~3–4 h wall-clock with 5 teams parallel.

---

## 6. Phase 2 — Integration & Hardening

**Owner:** 1 integration agent (Team Lead role), sequential.

### Entry gates

- [ ] All 5 Phase 1 branches pass CI against `main`
- [ ] Each branch has green contract tests imported from `contract-v1`
- [ ] No open changes touching `internal/contract/` or `internal/graphql/schema.graphql`

### Steps

1. **Merge** branches in order: `infra/core` → `read/queries` → `write/mutations` → `realtime/ws` → `release/infra`. Resolve interface drift against `docs/interface-contract.md`. Run `go generate ./...` after each merge.
2. **Fixture recording** — for every command, record one happy-path go-vcr cassette at `testdata/fixtures/<command>_happy.yaml` against a dedicated test Fireflies account. Fixtures scrubbed of `Bearer `, email addresses, and UUIDs using `scripts/scrub-fixtures.sh`. Cassettes committed; CI runs tests in replay-only mode.
3. **Full lint + test pass** — `golangci-lint run ./...`, `go vet ./...`, `go test -race ./...`, `govulncheck ./...` all exit 0.
4. **Security review** (see §9 for full list) — constant-time HMAC, no key in logs, config 0600, destructive-op confirm.
5. **Token-efficiency benchmark** (see §8) — run `benchmarks/token_efficiency/run.sh` once manually before tagging v1.0.0; commit `results.json` alongside the release. Target: mean default-mode ratio ≤ 0.40 vs. MCP baseline. Not a CI gate (runs once per release).
6. **Docs** — `README.md` (install, quickstart, auth, full command reference), `docs/command-reference.md` auto-generated via `cobra.GenMarkdownTree`, `docs/op-command-map.md` table mapping every GraphQL op to its command. (Man pages deferred — not a v1.0.0 requirement.)

### Exit gates

All Phase 1 gates hold + full command-level acceptance checklist (§7) complete for 100 % of commands + token benchmark passes + security review signed off + docs generated.

**Effort:** 6–8 agent-hours.

---

## 7. Command-level acceptance checklist

Tracked in `docs/command-checklist.md` (rows = commands, columns = items). Every command must pass all items before Phase 2 exit.

- [ ] `--help` names the underlying GraphQL op
- [ ] Unit test exercising happy path via go-vcr fixture
- [ ] Fixture at `testdata/fixtures/<command>_happy.yaml`, scrubbed of secrets
- [ ] Returns exit 3 on auth failure with actionable stderr
- [ ] Respects `--profile`, `--json`, `--output`, `--jq`, `--fields`
- [ ] Table output has stable documented column order
- [ ] Mutations: `--dry-run` prints op + vars, zero HTTP (asserted via cassette hit count 0)
- [ ] Destructive ops: confirm prompt; `--yes` bypasses
- [ ] Transcript-returning commands honour `--transcript`
- [ ] `go vet` clean

---

## 8. Contract tests & token-efficiency benchmark

### Contract tests (`internal/contract/contract_test.go`, must pass before any Phase 1 branch diverges)

- `TestClientAuth_EnvKey`, `TestClientAuth_ConfigFile` — auth header precedence
- `TestClientRetry_503` — 3 retries then success
- `TestClientRateLimit_429` — honours `Retry-After`; completes in < 2.5 s
- `TestClientFieldSelection` — `--fields id,title` produces a GraphQL doc containing only those fields (assert on recorded request body)
- `TestRenderer_{Table,JSON,NDJSON,YAML,TSV,Plaintext}` — valid output per format, stable column order, no ANSI in plaintext
- `TestConfig_Precedence`, `TestConfig_ProfileFlag`, `TestConfig_MissingKey`
- `TestExitCode_{AuthFailure,NotFound,RateLimit,ServerError,UserAbort}` — exit codes match §3
- `TestFlags_{JSONExclusive,FieldsCSV,TranscriptEnum}` — shared flag parsing identical across all commands

### Token-efficiency benchmark (`benchmarks/token_efficiency/`)

Corpus of 10 representative tasks (list-meetings, transcript-text, search, summary, analytics, action-items, users-list, fields-filtered get, soundbites-list, whoami). For each task and each mode (default table, `--json`, `--fields`, `--jq`), measure output bytes against a recorded MCP-response baseline at `mcp_baseline/*.json`.

**Pass criteria (manual gate at release time, not CI-enforced):**
- Default-mode mean ratio `cli_bytes / mcp_bytes` ≤ **0.40**
- `--fields` mean ratio ≤ **0.20**

This benchmark is the explicit justification for building the CLI — the whole project exists to reduce tokens vs. the MCP path. It is run once before every tagged release and the `results.json` is attached to the GitHub release. If a release misses the target, the regression is a ship blocker.

---

## 9. Risk register

| # | Risk | Sev | Lkhd | Mitigation |
|---|---|---|---|---|
| 1 | Fireflies API undocumented / breaking field changes | High | Med | Pinned `schema.graphql`; weekly `schema-drift-check` CI job diffs live introspection, opens issue on delta |
| 2 | Rate-limit lockout during CI | High | High | CI tests are go-vcr replay only; live-API tests gated behind `TEST_LIVE=1`, never set in CI |
| 3 | Webhook signature bypass (security-critical) | Crit | Low | `hmac.Equal` mandatory; `TestWebhookSignatureTimingConstant`; grep for `== sig` in CI fails build |
| 4 | API key / PII in committed fixtures | Crit | Med | `scripts/scrub-fixtures.sh --check` in CI; greps for `Bearer `, `@`, UUID, fails on hit |
| 5 | Socket.IO Go lib maturity | Med | Med | Isolate behind `internal/realtime` interface; realtime integration test skipped in CI |
| 6 | Binary size / cold-start regression | Med | Low | CI fails if any target > 20 MB or median cold start > 50 ms on macOS arm64 |
| 7 | Homebrew formula publish failure | Med | Med | `scripts/test-homebrew.sh` runs `brew install` on macOS runner before tagging |
| 8 | Interface contract drift between teams | High | Med | `internal/contract/` branch-protected, import-only after `contract-v1`; RFC required for changes |
| 9 | genqlient schema drift during dev | High | Med | `schema-drift-check` job hard-fails if live introspection hash ≠ committed schema without matching commit |
| 10 | Shell completion breakage | Low | Low | CI pipes `fireflies completion zsh` to `zsh --no-exec` for syntax check |

---

## 10. Security gates (all must pass before `v1.0.0`)

- [ ] `scripts/scrub-fixtures.sh --check` exits 0 on `testdata/fixtures/`
- [ ] `TestWebhookSignatureHMAC` asserts `hmac.Equal`; `grep -r '== sig' internal/webhook/` finds nothing
- [ ] `TestAuthTokenNotLogged` — `FIREFLIES_API_KEY=secret123 fireflies --verbose meetings list` contains no `secret123` in stdout/stderr
- [ ] `TestConfigFilePermissions` — config written with mode `0600`
- [ ] `govulncheck ./...` exits 0 on release branch

---

## 11. Phase 3 — Release `v1.0.0`

**Owner:** Team Release agent, sequential.

1. `make release-dry` — goreleaser snapshot builds all 5 targets; SBOM + changelog render.
2. Confirm CI green on `main` and `docs/op-command-map.md` is 34/34.
3. `git tag v1.0.0 && git push origin v1.0.0` → triggers release workflow.
4. goreleaser publishes archives + `checksums.txt` + `fireflies_1.0.0_sbom.json` to GitHub Releases.
5. Update `Formula/fireflies.rb` in `homebrew-tap` with `v1.0.0` URL + SHA256 from published `checksums.txt`; commit + push.
6. `brew install fvdm-otinga/tap/fireflies` smoke test on a fresh macOS arm64 runner.
7. README install section updated with `go install github.com/fvdm-otinga/fireflies-cli@latest` and `brew install …` lines.

### v1.0.0 ship criteria

- [ ] All 34 GraphQL ops covered (`docs/op-command-map.md` 34/34)
- [ ] Token-efficiency benchmark passes (≤ 0.40 default, ≤ 0.20 `--fields`)
- [ ] `golangci-lint` + `go vet` + `govulncheck` clean
- [ ] CI matrix green on `{linux/amd64, linux/arm64, darwin/amd64, darwin/arm64}`; windows builds compile
- [ ] `README.md` complete; command reference auto-built
- [ ] `brew install fvdm-otinga/tap/fireflies` succeeds on a clean macOS arm64 runner
- [ ] Binary ≤ 20 MB on all targets; cold start ≤ 50 ms on macOS arm64
- [ ] All security gates pass
- [ ] Command-level acceptance checklist 100 % complete
- [ ] `CHANGELOG.md` entry for v1.0.0 listing all commands

**Effort:** 2–3 agent-hours.

---

## 12. Summary

| Phase | Wall time | Agent-hours | Teams |
|---|---|---|---|
| 0 — Contract & scaffold | Sequential | 3–4 | 1 architect |
| 1 — Parallel build | ~3–4 h | 12–16 | 5 teams (8 agents) |
| 2 — Integration & hardening | Sequential | 6–8 | 1 integrator |
| 3 — Release | Sequential | 2–3 | 1 release agent |
| **Total** | | **23–31** | |

**Ships:** ~45 named subcommands covering all 17 queries + 17 mutations + realtime (`realtime tail`, `live add/soundbite/action-item`) + webhooks (`webhooks serve`) + housekeeping (`auth {login,logout,status}`, `config {get,set,list}`, `version`, `completion`).

---

## 13. Assumptions to verify at execution time

These were assumed during planning and deserve a sanity check when Phase 0 starts:

- A valid Fireflies Business-tier API key is available for schema introspection and for recording go-vcr fixtures in Phase 2. If only Pro-tier access is available, `analytics` and `rule_executions_by_meeting` will be gated and must be tagged as "plan-gated" in help output.
- A recorded MCP-baseline corpus exists for the token-efficiency benchmark. If not, generate it in Phase 0 by running the 10 corpus tasks through the current Fireflies MCP tools and capturing raw JSON responses to `benchmarks/token_efficiency/mcp_baseline/*.json`.
- GitHub personal token for `fvdm-otinga` has `repo`, `workflow`, and `write:packages` scopes (already confirmed: current token has `repo, workflow, read:org, gist`; `write:packages` needs adding before Phase 3).
- `github.com/fvdm-otinga/homebrew-tap` does not yet exist and must be created in Phase 0 (empty repo with an `Formula/` directory).
- The Fireflies Socket.IO realtime API is accessible without a separate feature flag on the account.
- No account-level rate limit exists beyond the documented per-op limits. If one is found, the `--profile` flag enables swapping to a second account for CI recordings.
