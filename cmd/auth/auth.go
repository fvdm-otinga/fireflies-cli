// Package auth implements the `fireflies auth` command group.
package auth

import "github.com/spf13/cobra"

// NewAuthCmd returns the `auth` command group.
func NewAuthCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: "Manage Fireflies API authentication",
	}
	c.AddCommand(newLoginCmd())
	c.AddCommand(newLogoutCmd())
	c.AddCommand(newStatusCmd())
	return c
}
