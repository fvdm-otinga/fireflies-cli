// Package soundbites implements the `fireflies soundbites ...` command group.
package soundbites

import "github.com/spf13/cobra"

// NewSoundbitesCmd returns the `soundbites` command group.
func NewSoundbitesCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "soundbites",
		Short: "Query Fireflies soundbites (bites)",
	}
	c.AddCommand(newListCmd())
	c.AddCommand(newGetCmd())
	return c
}
