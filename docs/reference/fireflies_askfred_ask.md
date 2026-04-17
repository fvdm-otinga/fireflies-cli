## fireflies askfred ask

Ask a question about a meeting via AskFred (GraphQL: createAskFredThread)

### Synopsis

Ask AskFred a question about a meeting transcript.

The question can be provided as a positional argument or piped via stdin:
  fireflies askfred ask --meeting <id> "summarize in one sentence"
  echo "summarize in one sentence" | fireflies askfred ask --meeting <id>

```
fireflies askfred ask [question] [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for ask
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --meeting string      Meeting (transcript) ID to ask about (required)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies askfred](fireflies_askfred.md)	 - Query AskFred AI threads

