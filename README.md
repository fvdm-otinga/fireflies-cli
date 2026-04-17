# fireflies

A token-efficient Go CLI for the [Fireflies.ai](https://fireflies.ai) GraphQL API.

Built for use from LLM / agent workflows (Claude Code, Cursor, scripts). Default output is human-readable; `--json` / `--jq` / `--fields` enable machine-friendly, compact responses.

## Status

Pre-release. See [`plan-fireflies-cli.md`](./plan-fireflies-cli.md) for the execution plan.

## Install

```sh
# Go toolchain
go install github.com/fvdm-otinga/fireflies-cli@latest

# Homebrew (tap)
brew install fvdm-otinga/tap/fireflies
```

Prebuilt binaries for macOS (arm64 / amd64), Linux (amd64 / arm64), and Windows (amd64) are published on each [GitHub Release](https://github.com/fvdm-otinga/fireflies-cli/releases).

## Quick start

```sh
# 1. Authenticate — key is generated at https://app.fireflies.ai/integrations/custom/fireflies
export FIREFLIES_API_KEY=ff_...

# 2. Verify
fireflies users whoami
```

Or use the interactive login:

```sh
fireflies auth login
```

## Shared flags

Every command accepts:

| Flag | Purpose |
|---|---|
| `--profile` | Config profile name (multi-account support) |
| `--json` | Shortcut for `--output json` |
| `--output table\|json\|ndjson\|yaml\|tsv\|plaintext` | Output format (default `table`) |
| `--jq <expr>` | Post-process output via embedded `gojq` |
| `--fields a,b,c` | Keep only the named fields (client-side projection) |
| `--limit`, `--skip` | Pagination |
| `--transcript none\|preview\|full` | Transcript depth (list→none, get→preview) |
| `--since`, `--until` | Time-range filters (RFC3339 or `7d`) |
| `--dry-run` | (mutations) print the GraphQL op without executing |
| `--yes` | (destructive) bypass confirmation prompt |

## Documentation

- Command reference — `docs/command-reference.md` (auto-generated)
- Interface contract for contributors — `docs/interface-contract.md`
- Op-to-command map — `docs/op-command-map.md`

## License

MIT — see [LICENSE](./LICENSE).
