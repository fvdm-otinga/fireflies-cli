## fireflies

Fireflies.ai CLI (token-efficient wrapper for the GraphQL API)

### Synopsis

fireflies is a command-line interface for the Fireflies.ai GraphQL API,
designed for efficient use from LLM/agent workflows.

Default output is a human-readable table; use --json for machine output.
All commands accept --fields (field projection), --jq (post-filter),
--output (format), and --profile (config profile).

### Options

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
  -h, --help                help for fireflies
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

* [fireflies analytics](fireflies_analytics.md)	 - Fetch workspace analytics (GraphQL: analytics)
* [fireflies apps](fireflies_apps.md)	 - Query Fireflies app outputs
* [fireflies askfred](fireflies_askfred.md)	 - Query AskFred AI threads
* [fireflies auth](fireflies_auth.md)	 - Manage Fireflies API authentication
* [fireflies channels](fireflies_channels.md)	 - Query Fireflies channels
* [fireflies completion](fireflies_completion.md)	 - Generate shell completion scripts
* [fireflies config](fireflies_config.md)	 - Manage Fireflies CLI configuration
* [fireflies contacts](fireflies_contacts.md)	 - Query Fireflies contacts
* [fireflies live](fireflies_live.md)	 - Query and interact with live meetings
* [fireflies meetings](fireflies_meetings.md)	 - Query and manage Fireflies meetings (transcripts)
* [fireflies realtime](fireflies_realtime.md)	 - Stream live transcript events via Socket.IO
* [fireflies rules](fireflies_rules.md)	 - Query Fireflies automation rule executions
* [fireflies soundbites](fireflies_soundbites.md)	 - Query Fireflies soundbites (bites)
* [fireflies transcript](fireflies_transcript.md)	 - Work with meeting transcripts
* [fireflies users](fireflies_users.md)	 - Manage and query Fireflies users
* [fireflies version](fireflies_version.md)	 - Print the CLI version, commit, and build date
* [fireflies webhooks](fireflies_webhooks.md)	 - Webhook utilities (receive and verify Fireflies webhook events)

