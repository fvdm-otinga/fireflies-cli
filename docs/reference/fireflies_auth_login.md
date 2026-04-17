## fireflies auth login

Save a Fireflies API key to the config file

### Synopsis

Prompt for a Fireflies API key (input is not echoed), verify it via
the Whoami query, and write it to ~/.config/fireflies/config.toml with
file mode 0600 under the given profile (default: "default").

```
fireflies auth login [flags]
```

### Options

```
  -h, --help             help for login
      --profile string   Profile name to save the key under (default: "default")
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
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies auth](fireflies_auth.md)	 - Manage Fireflies API authentication

