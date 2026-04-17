// Package rules implements the `fireflies rules ...` command group.
package rules

import "github.com/spf13/cobra"

// NewRulesCmd returns the `rules` command group.
func NewRulesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "rules",
		Short: "Query Fireflies automation rule executions",
	}
	c.AddCommand(newExecutionsCmd())
	return c
}
