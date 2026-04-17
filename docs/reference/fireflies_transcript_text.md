## fireflies transcript text

Print meeting transcript as plaintext (GraphQL: transcript.sentences)

### Synopsis

Fetch and print the full transcript text for a meeting in Speaker: text format.

  --format plaintext  (default) — one line per sentence: "Speaker: text"
  --format json       — raw sentences JSON array
  --since / --until   — filter sentences by start_time (RFC3339 window)

```
fireflies transcript text <id> [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for text
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

* [fireflies transcript](fireflies_transcript.md)	 - Work with meeting transcripts

