// Package cmd implements the `fireflies` CLI root command.
package cmd

import (
	"github.com/spf13/cobra"

	usercmd "github.com/fvdm-otinga/fireflies-cli/cmd/users"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
)

// NewRootCmd returns the root `fireflies` Cobra command.
func NewRootCmd(version string) *cobra.Command {
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

	root.AddCommand(usercmd.NewUsersCmd())
	root.AddCommand(newVersionCmd(version))
	return root
}

func newVersionCmd(version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Printf("fireflies version %s\n", version)
		},
	}
}
