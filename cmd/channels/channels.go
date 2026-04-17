// Package channels implements the `fireflies channels ...` command group.
package channels

import "github.com/spf13/cobra"

// NewChannelsCmd returns the `channels` command group.
func NewChannelsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "channels",
		Short: "Query Fireflies channels",
	}
	c.AddCommand(newListCmd())
	c.AddCommand(newGetCmd())
	return c
}
