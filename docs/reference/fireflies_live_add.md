## fireflies live add

Add a participant to a live meeting (GraphQL: addToLiveMeeting, rate: 3/20min)

### Synopsis

Adds an attendee to a live/active Fireflies meeting.

The --meeting flag accepts the meeting link (URL) that Fireflies uses to
identify the live session, not the transcript ID.

Rate limit: 3 requests per 20 minutes (enforced client-side).

```
fireflies live add [flags]
```

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --email string        Email of the attendee to add
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for add
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --meeting string      Meeting link (URL) of the live meeting (required)
      --name string         Display name of the attendee to add
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies live](fireflies_live.md)	 - Query and interact with live meetings

