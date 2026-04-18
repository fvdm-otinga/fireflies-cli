# fireflies

A token-efficient Go CLI for the [Fireflies.ai](https://fireflies.ai) GraphQL API.

Built for LLM / agent workflows (Claude Code, Cursor, scripts). Default output is a
human-readable table; `--json` / `--jq` / `--fields` give compact, machine-friendly
responses that use **60–80% fewer tokens** than the equivalent MCP tool calls.

---

## Install

### One-liner (interactive wizard, macOS / Linux)

```sh
curl -fsSL https://raw.githubusercontent.com/fvdm-otinga/fireflies-cli/main/scripts/install.sh | bash
```

Detects OS/arch, verifies the SHA-256 checksum, installs to `~/.local/bin` (or a
prefix you pick), and optionally runs `fireflies auth login` plus shell
completions. Non-interactive mode: `… | bash -s -- --yes`.

### Homebrew (macOS / Linux)

```sh
brew install --cask fvdm-otinga/tap/fireflies
```

### Go install

```sh
go install github.com/fvdm-otinga/fireflies-cli@latest
```

### Prebuilt binary

Download from [GitHub Releases](https://github.com/fvdm-otinga/fireflies-cli/releases)
for macOS (arm64 / amd64), Linux (amd64 / arm64), and Windows (amd64).

See [docs/install.md](./docs/install.md) for detailed instructions and checksum verification.

---

## Auth

```sh
# Option A — environment variable (CI / scripts)
export FIREFLIES_API_KEY=ff_xxxxxxxxxxxxxxxx
fireflies users whoami

# Option B — interactive login (persisted in ~/.config/fireflies/config.toml)
fireflies auth login
fireflies auth status

# Multi-profile
fireflies auth login --profile work
fireflies meetings list --profile work
```

Generate your API key at [app.fireflies.ai/integrations/custom/fireflies](https://app.fireflies.ai/integrations/custom/fireflies).

---

## Quickstart

```sh
# Who am I?
fireflies users whoami

# List recent meetings
fireflies meetings list

# Get a meeting as JSON
fireflies meetings get <meeting-id> --json

# Full transcript as speaker-attributed plaintext
fireflies transcript text <meeting-id>

# Filter output with jq
fireflies meetings list --json --jq '.[] | {id, title, date}'

# Keep only specific fields (saves tokens)
fireflies meetings list --fields id,title,date

# Watch live/active meetings
fireflies meetings active

# Upload a recording
fireflies meetings upload ./recording.mp3
fireflies meetings upload https://example.com/meeting.mp4

# Stream realtime transcript events
fireflies realtime tail <meeting-id>

# Receive webhooks locally (secret read from FIREFLIES_WEBHOOK_SECRET)
FIREFLIES_WEBHOOK_SECRET="$WEBHOOK_SECRET" fireflies webhooks serve --port 8080
```

---

## Commands

```
fireflies
├── auth           login / logout / status
├── config         get / set / list
├── users          whoami / list / groups / set-role
├── meetings       list / get / active / upload / update / move / share / revoke / delete
├── channels       list / get
├── transcript     text <id>    ← plaintext speaker-attributed transcript
├── analytics      (workspace analytics)
├── soundbites     list / get / create
├── apps           list
├── askfred        threads / thread / ask / continue / delete
├── rules          executions <meeting-id>
├── contacts       list
├── live           items / add / soundbite / action-item
├── realtime       tail <meeting-id>
├── webhooks       serve
├── completion     bash / zsh / fish / powershell
└── version
```

For the full per-command reference, see:

- [docs/command-reference.md](./docs/command-reference.md) — annotated command tree
- [docs/reference/](./docs/reference/) — auto-generated per-command Markdown (run `make docs`)
- [docs/op-command-map.md](./docs/op-command-map.md) — all 35 GraphQL ops → CLI commands

---

## Shared flags

Every command accepts:

| Flag | Description |
|---|---|
| `--profile` | Config profile (multi-account) |
| `--json` | Shortcut for `--output json` |
| `--output table\|json\|ndjson\|yaml\|tsv\|plaintext` | Output format |
| `--jq <expr>` | Post-process via embedded gojq |
| `--fields a,b,c` | Client-side field projection (reduces token count) |
| `--limit`, `--skip` | Pagination |
| `--transcript none\|preview\|full` | Transcript verbosity (list→none, get→preview) |
| `--since`, `--until` | Time-range filters (RFC3339 or `7d`) |
| `--dry-run` | Print GraphQL op without executing (mutations) |
| `--yes` | Bypass confirmation prompts (destructive ops) |

---

## Token-efficiency rationale

Fireflies MCP tools return full JSON payloads including every transcript sentence,
summary section, and analytics field. For an agent querying "list my last 5 meetings"
the raw MCP response can be 20–50 KB. The CLI default table output for the same query
is typically 500 bytes.

Key mechanisms:
- **Default table output** — human-readable, no transcript sentences
- **`--transcript none|preview|full`** — controls whether and how much transcript
  text is fetched (default on `list` is `none`)
- **`--fields a,b,c`** — client-side projection removes unused fields before rendering
- **`--jq <expr>`** — embedded gojq expression filters the JSON in-process, no subprocess

Target: CLI default ≤ 0.40× MCP baseline tokens; `--fields` mode ≤ 0.20×.

---

## Contributing

See [docs/interface-contract.md](./docs/interface-contract.md) for the architecture contract.

```sh
git clone https://github.com/fvdm-otinga/fireflies-cli.git
cd fireflies-cli
make build   # produces ./fireflies
make test    # go test -race ./...
make lint    # golangci-lint run ./...
make docs    # regenerates docs/reference/
```

Do not modify `internal/contract/` or `internal/graphql/schema.graphql` without an RFC
in `docs/rfc/`.

---

## License

MIT — see [LICENSE](./LICENSE).
