// Package live implements the `fireflies live ...` command group.
package live

import "github.com/spf13/cobra"

// NewLiveCmd returns the `live` command group.
func NewLiveCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "live",
		Short: "Query and interact with live meetings",
	}
	c.AddCommand(newItemsCmd())
	c.AddCommand(newAddCmd())
	c.AddCommand(newSoundbiteCmd())
	c.AddCommand(newActionItemCmd())
	return c
}
