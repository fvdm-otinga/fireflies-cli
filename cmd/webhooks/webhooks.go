// Package webhooks implements the `fireflies webhooks` command group.
package webhooks

import "github.com/spf13/cobra"

// NewWebhooksCmd returns the `webhooks` command group.
func NewWebhooksCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "webhooks",
		Short: "Webhook utilities (receive and verify Fireflies webhook events)",
	}
	c.AddCommand(newServeCmd())
	return c
}
