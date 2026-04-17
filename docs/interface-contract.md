# Interface Contract (frozen at `contract-v1`)

Every command in this CLI builds against the interfaces in `internal/contract/`. These interfaces are **frozen** at the git tag `contract-v1`. Changes require a documented RFC in `docs/rfc/` and approval from the repo owner.

## Why a contract

The CLI is built by five parallel teams (§ `plan-fireflies-cli.md`). The contract is the shared interface they all consume — it exists so Team Read, Team Write, Team Realtime, Team Release, and Team Infra can ship independently without integration drift.

## Shared packages

| Package | Purpose |
|---|---|
| `internal/client`    | GraphQL HTTP client with auth, retry, per-op rate buckets |
| `internal/config`    | TOML profile loader; env `FIREFLIES_API_KEY` overrides |
| `internal/output`    | Format renderer (table/json/ndjson/yaml/tsv/plaintext) + `--jq` + client-side projection |
| `internal/errors`    | Typed `CLIError` + exit-code taxonomy |
| `internal/flags`     | Shared Cobra flag binder (one call sets up every --* flag) |
| `internal/graphql`   | `genqlient` config and generated types |
| `internal/graphql/dynamic` | Hand-written dynamic selection-set builders for fat-payload ops |

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | General / API error (structured JSON on stderr) |
| 2 | Usage error (bad flag, invalid enum) |
| 3 | Auth error (missing/invalid API key) |
| 4 | Rate-limit exhausted |
| 5 | Not found |

Structured error shape on stderr: `{"code": "...", "error": "..."}`.

## `--fields` strategy

`genqlient` produces fixed typed query structs; it cannot narrow a query's selection set at runtime. Two paths:

1. **Default (client-side projection)** — applied by `internal/output.Project`. Keeps only the named top-level (and dotted-nested) paths after decode. Saves output tokens, not API bytes. Wired into every command automatically via the shared flags.
2. **Fat-payload path** — only for `meetings list`, `meetings get`, and `transcript text`. These commands import `internal/graphql/dynamic` to build the GraphQL document string at runtime, skipping heavy fields (`sentences`, `summary.*`, `analytics`) unless explicitly requested. Typed response-side structs are via a shared `Transcript` struct with pointer fields (nil = not requested).

Discovery: `--fields ?` prints the whitelist for that command.

## Renderer input shape

Commands always pass their **top-level `genqlient` response struct** (or a slice) to `output.Render`. The renderer normalises via `json.Marshal` → `json.Unmarshal([]byte, &any)` before projection or `gojq`. No command may build a map or custom shape and pass it to the renderer.

List responses may be wrapped in `{"data": [...], "meta": {"limit": N, "skip": N, "next_skip": N|null}}` via `internal/output.Envelope` (to be added in Phase 1). Single-object responses are emitted unwrapped.

## Column definitions

Table column order and header text for each response type live under `internal/output/columns/` (one file per type). Every command of that type imports the same column definitions — never inlines them.

## Shared flags (bound by `flags.Bind(cmd)`)

```
--profile string    --json          --jq string      --output string
--fields string     --limit int     --skip int
--transcript none|preview|full      --since RFC3339  --until RFC3339
--yes               --dry-run
```

## Changing the contract

1. Write an RFC at `docs/rfc/NNNN-short-name.md` explaining the motivation, the proposed change, impact on each team, and migration plan.
2. Open a PR changing `internal/contract/` with the RFC link in the description.
3. Require owner approval (CODEOWNERS). Bump the tag to `contract-v2` on merge.
