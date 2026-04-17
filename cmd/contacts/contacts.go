// Package contacts implements the `fireflies contacts ...` command group.
package contacts

import "github.com/spf13/cobra"

// NewContactsCmd returns the `contacts` command group.
func NewContactsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "contacts",
		Short: "Query Fireflies contacts",
	}
	c.AddCommand(newListCmd())
	return c
}
