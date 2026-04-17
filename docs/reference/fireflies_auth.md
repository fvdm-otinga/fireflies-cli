## fireflies auth

Manage Fireflies API authentication

### Options

```
  -h, --help   help for auth
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
* [fireflies auth login](fireflies_auth_login.md)	 - Save a Fireflies API key to the config file
* [fireflies auth logout](fireflies_auth_logout.md)	 - Remove a profile from the config file
* [fireflies auth status](fireflies_auth_status.md)	 - Print the current authentication status

