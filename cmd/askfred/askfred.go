// Package askfred implements the `fireflies askfred ...` command group.
package askfred

import "github.com/spf13/cobra"

// NewAskFredCmd returns the `askfred` command group.
func NewAskFredCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "askfred",
		Short: "Query AskFred AI threads",
	}
	c.AddCommand(newThreadsCmd())
	c.AddCommand(newThreadCmd())
	c.AddCommand(newAskCmd())
	c.AddCommand(newContinueCmd())
	c.AddCommand(newDeleteCmd())
	return c
}
