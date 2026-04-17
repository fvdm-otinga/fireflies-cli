## fireflies completion

Generate shell completion scripts

### Synopsis

Generate shell completion scripts for fireflies.

Bash:
  # Add to ~/.bash_profile or ~/.bashrc:
  source <(fireflies completion bash)

  # Or write to a file:
  fireflies completion bash > /usr/local/etc/bash_completion.d/fireflies

Zsh:
  # Enable completions if not already done (add to ~/.zshrc):
  autoload -Uz compinit && compinit

  # Add to ~/.zshrc:
  source <(fireflies completion zsh)

  # Or install to your completions directory:
  fireflies completion zsh > "${fpath[1]}/_fireflies"

Fish:
  fireflies completion fish | source

  # Or persist:
  fireflies completion fish > ~/.config/fish/completions/fireflies.fish

PowerShell:
  fireflies completion powershell | Out-String | Invoke-Expression

  # Or write to profile:
  fireflies completion powershell > fireflies.ps1
  # then dot-source it from $PROFILE


```
fireflies completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
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

