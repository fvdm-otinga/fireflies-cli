## fireflies auth status

Print the current authentication status

### Synopsis

Print the active profile name, config file path, whether FIREFLIES_API_KEY
is set (masked), and the authenticated user's email (via Whoami).

```
fireflies auth status [flags]
```

### Options

```
  -h, --help   help for status
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

* [fireflies auth](fireflies_auth.md)	 - Manage Fireflies API authentication

