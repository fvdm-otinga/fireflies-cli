// Package cmd implements the `fireflies` CLI root command.
package cmd

import (
	"github.com/spf13/cobra"

	authcmd "github.com/fvdm-otinga/fireflies-cli/cmd/auth"
	cfgcmd "github.com/fvdm-otinga/fireflies-cli/cmd/config"
	usercmd "github.com/fvdm-otinga/fireflies-cli/cmd/users"
	vercmd "github.com/fvdm-otinga/fireflies-cli/cmd/version"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
)

// NewRootCmd returns the root `fireflies` Cobra command.
func NewRootCmd(version, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:   "fireflies",
		Short: "Fireflies.ai CLI (token-efficient wrapper for the GraphQL API)",
		Long: `fireflies is a command-line interface for the Fireflies.ai GraphQL API,
designed for efficient use from LLM/agent workflows.

Default output is a human-readable table; use --json for machine output.
All commands accept --fields (field projection), --jq (post-filter),
--output (format), and --profile (config profile).`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	flags.Bind(root)
	root.Version = version
	root.SetVersionTemplate("fireflies version {{.Version}}\n")

	root.AddCommand(authcmd.NewAuthCmd())
	root.AddCommand(cfgcmd.NewConfigCmd())
	root.AddCommand(usercmd.NewUsersCmd())
	root.AddCommand(vercmd.NewVersionCmd(version, commit, date))
	return root
}
