## fireflies live soundbite

Create a live soundbite for a meeting (GraphQL: createLiveSoundbite)

### Synopsis

Creates a live soundbite for an active meeting.

The --prompt flag describes the soundbite content. Use --start and --end
(seconds) to specify a time window; these are folded into the prompt when
--prompt is not set.

```
fireflies live soundbite [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --end int             End time in seconds
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for soundbite
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --meeting string      Meeting ID (required)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --prompt string       Soundbite prompt/description
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --start int           Start time in seconds (combined with --end if --prompt not set)
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies live](fireflies_live.md)	 - Query and interact with live meetings

