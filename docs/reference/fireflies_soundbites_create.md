## fireflies soundbites create

Create a soundbite from a meeting (GraphQL: createBite)

### Synopsis

Create a soundbite clip from a meeting transcript.

Time values (--start, --end) accept:
  - Millisecond epoch:   30000
  - mm:ss format:        0:30 or 1:15

```
fireflies soundbites create [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --end string          End time in ms or mm:ss (required)
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for create
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --media-type string   Media type (audio or video)
      --meeting string      Meeting (transcript) ID (required)
      --name string         Name for the soundbite
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --start string        Start time in ms or mm:ss (required)
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies soundbites](fireflies_soundbites.md)	 - Query Fireflies soundbites (bites)

