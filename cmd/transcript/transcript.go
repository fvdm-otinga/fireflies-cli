// Package transcript implements the `fireflies transcript ...` command group.
package transcript

import "github.com/spf13/cobra"

// NewTranscriptCmd returns the `transcript` command group.
func NewTranscriptCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "transcript",
		Short: "Work with meeting transcripts",
	}
	c.AddCommand(newTextCmd())
	return c
}
