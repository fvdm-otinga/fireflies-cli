// Package realtime implements the `fireflies realtime` command group.
package realtime

import "github.com/spf13/cobra"

// NewRealtimeCmd returns the `realtime` command group.
func NewRealtimeCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "realtime",
		Short: "Stream live transcript events via Socket.IO",
	}
	c.AddCommand(newTailCmd())
	return c
}
