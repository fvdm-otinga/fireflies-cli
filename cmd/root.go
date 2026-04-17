// Package cmd implements the `fireflies` CLI root command.
package cmd

import (
	"github.com/spf13/cobra"

	analytcmd "github.com/fvdm-otinga/fireflies-cli/cmd/analytics"
	appscmd "github.com/fvdm-otinga/fireflies-cli/cmd/apps"
	askfredcmd "github.com/fvdm-otinga/fireflies-cli/cmd/askfred"
	authcmd "github.com/fvdm-otinga/fireflies-cli/cmd/auth"
	cfgcmd "github.com/fvdm-otinga/fireflies-cli/cmd/config"
	channelscmd "github.com/fvdm-otinga/fireflies-cli/cmd/channels"
	contactscmd "github.com/fvdm-otinga/fireflies-cli/cmd/contacts"
	livecmd "github.com/fvdm-otinga/fireflies-cli/cmd/live"
	meetingscmd "github.com/fvdm-otinga/fireflies-cli/cmd/meetings"
	realtimecmd "github.com/fvdm-otinga/fireflies-cli/cmd/realtime"
	rulescmd "github.com/fvdm-otinga/fireflies-cli/cmd/rules"
	soundbitescmd "github.com/fvdm-otinga/fireflies-cli/cmd/soundbites"
	transcriptcmd "github.com/fvdm-otinga/fireflies-cli/cmd/transcript"
	usercmd "github.com/fvdm-otinga/fireflies-cli/cmd/users"
	vercmd "github.com/fvdm-otinga/fireflies-cli/cmd/version"
	webhookscmd "github.com/fvdm-otinga/fireflies-cli/cmd/webhooks"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
)

// NewRootCmd returns the root `fireflies` Cobra command.
func NewRootCmd(version, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:   "fireflies",
		Short: "Fireflies.ai CLI (token-efficient wrapper for the GraphQL API)",
		Long: `fireflies is a command-line interface for the Fireflies.ai GraphQL API,
designed for efficient use from LLM/agent workflows.

Default output is a human-readable table; use --json for machine output.
All commands accept --fields (field projection), --jq (post-filter),
--output (format), and --profile (config profile).`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	flags.Bind(root)
	root.Version = version
	root.SetVersionTemplate("fireflies version {{.Version}}\n")

	root.AddCommand(authcmd.NewAuthCmd())
	root.AddCommand(cfgcmd.NewConfigCmd())
	root.AddCommand(usercmd.NewUsersCmd())
	root.AddCommand(vercmd.NewVersionCmd(version, commit, date))
	root.AddCommand(meetingscmd.NewMeetingsCmd())
	root.AddCommand(channelscmd.NewChannelsCmd())
	root.AddCommand(analytcmd.NewAnalyticsCmd())
	root.AddCommand(transcriptcmd.NewTranscriptCmd())
	root.AddCommand(soundbitescmd.NewSoundbitesCmd())
	root.AddCommand(appscmd.NewAppsCmd())
	root.AddCommand(askfredcmd.NewAskFredCmd())
	root.AddCommand(rulescmd.NewRulesCmd())
	root.AddCommand(contactscmd.NewContactsCmd())
	root.AddCommand(livecmd.NewLiveCmd())
	root.AddCommand(realtimecmd.NewRealtimeCmd())
	root.AddCommand(webhookscmd.NewWebhooksCmd())
	return root
}
