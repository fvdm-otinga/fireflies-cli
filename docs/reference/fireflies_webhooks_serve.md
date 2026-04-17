## fireflies webhooks serve

Start an HTTP server that receives and verifies Fireflies webhook events (V1 and V2)

### Synopsis

Starts an HTTP server that listens for Fireflies webhook POST requests.

Verified events are emitted as NDJSON on stdout.
Rejected events (bad or missing signature) return 401 and log a warning to stderr.

Routes:
  POST /webhooks/v1  — Fireflies webhook V1
  POST /webhooks/v2  — Fireflies webhook V2
  GET  /health       — Health check (200 ok)

The webhook secret must be set via --secret-env (env var name) or --secret.

```
fireflies webhooks serve [flags]
```

### Options

```
  -h, --help                help for serve
      --port int            Port to listen on (default 8080)
      --secret string       Webhook secret value (use --secret-env in production)
      --secret-env string   Name of env var holding the webhook secret
```

### Options inherited from parent commands

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies webhooks](fireflies_webhooks.md)	 - Webhook utilities (receive and verify Fireflies webhook events)

