## fireflies auth logout

Remove a profile from the config file

### Synopsis

Remove the named profile (default: active profile) from ~/.config/fireflies/config.toml. Use --yes to skip the confirmation prompt.

```
fireflies auth logout [flags]
```

### Options

```
  -h, --help             help for logout
      --profile string   Profile to remove (default: active profile)
      --yes              Bypass confirmation prompt
```

### Options inherited from parent commands

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
```

### SEE ALSO

* [fireflies auth](fireflies_auth.md)	 - Manage Fireflies API authentication

