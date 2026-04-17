## fireflies meetings update

Update meeting metadata

### Options

```
  -h, --help   help for update
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

* [fireflies meetings](fireflies_meetings.md)	 - Query and manage Fireflies meetings (transcripts)
* [fireflies meetings update privacy](fireflies_meetings_update_privacy.md)	 - Update meeting privacy level (GraphQL: updateMeetingPrivacy)
* [fireflies meetings update state](fireflies_meetings_update_state.md)	 - Pause or resume a live meeting recording (GraphQL: updateMeetingState)
* [fireflies meetings update title](fireflies_meetings_update_title.md)	 - Update meeting title (GraphQL: updateMeetingTitle)

