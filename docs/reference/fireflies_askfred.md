## fireflies askfred

Query AskFred AI threads

### Options

```
  -h, --help   help for askfred
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
* [fireflies askfred ask](fireflies_askfred_ask.md)	 - Ask a question about a meeting via AskFred (GraphQL: createAskFredThread)
* [fireflies askfred continue](fireflies_askfred_continue.md)	 - Continue an existing AskFred conversation thread (GraphQL: continueAskFredThread)
* [fireflies askfred delete](fireflies_askfred_delete.md)	 - Delete an AskFred thread and all its messages (GraphQL: deleteAskFredThread) [destructive]
* [fireflies askfred thread](fireflies_askfred_thread.md)	 - Get an AskFred thread by ID with all messages (GraphQL: askfred_thread)
* [fireflies askfred threads](fireflies_askfred_threads.md)	 - List AskFred conversation threads (GraphQL: askfred_threads)

