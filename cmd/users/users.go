// Package users implements the `fireflies users ...` command group.
package users

import "github.com/spf13/cobra"

// NewUsersCmd returns the `users` command group.
func NewUsersCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "users",
		Short: "Manage and query Fireflies users",
	}
	c.AddCommand(newWhoamiCmd())
	return c
}
