// Package config implements the `fireflies config` command group.
package config

import "github.com/spf13/cobra"

// NewConfigCmd returns the `config` command group.
func NewConfigCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Manage Fireflies CLI configuration",
	}
	c.AddCommand(newGetCmd())
	c.AddCommand(newSetCmd())
	c.AddCommand(newListCmd())
	return c
}
