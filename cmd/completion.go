// Package cmd provides the Cobra completion command for fireflies.
package cmd

import (
	"github.com/spf13/cobra"
)

// newCompletionCmd returns a completion command that generates shell completion
// scripts for bash, zsh, fish, and powershell. The completion subcommand is
// also auto-registered by Cobra; this explicit command provides customised
// long-form help text with installation instructions.
func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for fireflies.

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
`,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch args[0] {
			case "bash":
				return root.GenBashCompletionV2(cmd.OutOrStdout(), true)
			case "zsh":
				return root.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return root.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return root.GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			return nil
		},
	}
	return cmd
}
