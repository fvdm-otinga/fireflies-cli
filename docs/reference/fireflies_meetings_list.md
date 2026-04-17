## fireflies meetings list

List meetings/transcripts (GraphQL: transcripts)

### Synopsis

List meetings with pagination. The --transcript flag controls payload depth:
  none    — id, title, date, duration only (default)
  preview — adds organizer_email, participants
  full    — all fields including sentences, summary, analytics

Use --fields ? to list available dynamic fields.

```
fireflies meetings list [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for list
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

* [fireflies meetings](fireflies_meetings.md)	 - Query and manage Fireflies meetings (transcripts)

