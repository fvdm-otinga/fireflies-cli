## fireflies users

Manage and query Fireflies users

### Options

```
  -h, --help   help for users
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
* [fireflies users groups](fireflies_users_groups.md)	 - List user groups (GraphQL: user_groups)
* [fireflies users list](fireflies_users_list.md)	 - List all workspace users (GraphQL: users)
* [fireflies users set-role](fireflies_users_set-role.md)	 - Set a user's role (GraphQL: setUserRole) [destructive]
* [fireflies users whoami](fireflies_users_whoami.md)	 - Print the API key owner's profile (GraphQL: user)

