## fireflies live

Query and interact with live meetings

### Options

```
  -h, --help   help for live
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

* [fireflies](fireflies.md)	 - Fireflies.ai CLI (token-efficient wrapper for the GraphQL API)
* [fireflies live action-item](fireflies_live_action-item.md)	 - Create a live action item for a meeting (GraphQL: createLiveActionItem)
* [fireflies live add](fireflies_live_add.md)	 - Add a participant to a live meeting (GraphQL: addToLiveMeeting, rate: 3/20min)
* [fireflies live items](fireflies_live_items.md)	 - Get live action items for a meeting (GraphQL: live_action_items)
* [fireflies live soundbite](fireflies_live_soundbite.md)	 - Create a live soundbite for a meeting (GraphQL: createLiveSoundbite)

