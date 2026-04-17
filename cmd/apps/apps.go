// Package apps implements the `fireflies apps ...` command group.
package apps

import "github.com/spf13/cobra"

// NewAppsCmd returns the `apps` command group.
func NewAppsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "apps",
		Short: "Query Fireflies app outputs",
	}
	c.AddCommand(newListCmd())
	return c
}
