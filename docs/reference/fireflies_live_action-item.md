## fireflies live action-item

Create a live action item for a meeting (GraphQL: createLiveActionItem)

### Synopsis

Creates a live action item for an active meeting.

Use --prompt (or --text as an alias) for the action item description.
--assignee is recorded in the prompt if provided.

```
fireflies live action-item [flags]
```

### Options

```
      --assignee string     Assignee email (appended to prompt)
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for action-item
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --meeting string      Meeting ID (required)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --prompt string       Action item prompt/description
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --text string         Action item text (alias for --prompt)
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies live](fireflies_live.md)	 - Query and interact with live meetings

