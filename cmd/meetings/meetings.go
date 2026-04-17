// Package meetings implements the `fireflies meetings ...` command group.
package meetings

import "github.com/spf13/cobra"

// NewMeetingsCmd returns the `meetings` command group.
func NewMeetingsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "meetings",
		Short: "Query and manage Fireflies meetings (transcripts)",
	}
	c.AddCommand(newListCmd())
	c.AddCommand(newGetCmd())
	c.AddCommand(newActiveCmd())
	return c
}
