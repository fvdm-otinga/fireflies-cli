package meetings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newUpdateCmd returns the `fireflies meetings update` subcommand group.
func newUpdateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "update",
		Short: "Update meeting metadata",
	}
	c.AddCommand(newUpdateTitleCmd())
	c.AddCommand(newUpdatePrivacyCmd())
	c.AddCommand(newUpdateStateCmd())
	return c
}

// newUpdateTitleCmd returns `fireflies meetings update title <id> <title>`.
// GraphQL: UpdateMeetingTitle
func newUpdateTitleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "title <id> <title>",
		Short: "Update meeting title (GraphQL: updateMeetingTitle)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]
			title := args[1]

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{"id": id, "title": title},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation UpdateMeetingTitle($input: UpdateMeetingTitleInput!) {\n  updateMeetingTitle(input: $input) { id title date duration organizer_email }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.UpdateMeetingTitle(context.Background(), c, &ffgql.UpdateMeetingTitleInput{Id: id, Title: title})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				fmt.Fprintf(os.Stdout, "title updated: %s → %q\n", id, title)
				return nil
			}
			return output.Render(os.Stdout, resp.UpdateMeetingTitle, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "ID", Path: "id"},
					{Header: "TITLE", Path: "title"},
					{Header: "DATE", Path: "date"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}

// newUpdatePrivacyCmd returns `fireflies meetings update privacy <id> <level>`.
// GraphQL: UpdateMeetingPrivacy
func newUpdatePrivacyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "privacy <id> <link|owner|participants|teammatesandparticipants|teammates>",
		Short: "Update meeting privacy level (GraphQL: updateMeetingPrivacy)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]
			privacyStr := args[1]

			privacy := ffgql.MeetingPrivacy(privacyStr)
			switch privacy {
			case ffgql.MeetingPrivacyLink, ffgql.MeetingPrivacyOwner, ffgql.MeetingPrivacyParticipants,
				ffgql.MeetingPrivacyTeammatesandparticipants, ffgql.MeetingPrivacyTeammates:
				// valid
			default:
				return ferr.Usage(fmt.Sprintf("invalid privacy %q: must be one of link, owner, participants, teammatesandparticipants, teammates", privacyStr))
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{"id": id, "privacy": privacyStr},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation UpdateMeetingPrivacy($input: UpdateMeetingPrivacyInput!) {\n  updateMeetingPrivacy(input: $input) { id title privacy date }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.UpdateMeetingPrivacy(context.Background(), c, &ffgql.UpdateMeetingPrivacyInput{Id: id, Privacy: privacy})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				privVal := ""
				if resp.UpdateMeetingPrivacy.Privacy != nil {
					privVal = string(*resp.UpdateMeetingPrivacy.Privacy)
				}
				fmt.Fprintf(os.Stdout, "privacy updated: meeting %s → %s\n", id, privVal)
				return nil
			}
			return output.Render(os.Stdout, resp.UpdateMeetingPrivacy, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "ID", Path: "id"},
					{Header: "TITLE", Path: "title"},
					{Header: "PRIVACY", Path: "privacy"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}

// newUpdateStateCmd returns `fireflies meetings update state <id> <action>`.
// GraphQL: UpdateMeetingState
func newUpdateStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state <id> <pause_recording|resume_recording>",
		Short: "Pause or resume a live meeting recording (GraphQL: updateMeetingState)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			meetingID := args[0]
			actionStr := strings.ToLower(args[1])

			action := ffgql.MeetingStateAction(actionStr)
			switch action {
			case ffgql.MeetingStateActionPauseRecording, ffgql.MeetingStateActionResumeRecording:
				// valid
			default:
				return ferr.Usage(fmt.Sprintf("invalid action %q: must be pause_recording or resume_recording", actionStr))
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{"meeting_id": meetingID, "action": actionStr},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation UpdateMeetingState($input: UpdateMeetingStateInput!) {\n  updateMeetingState(input: $input) { success action }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.UpdateMeetingState(context.Background(), c, &ffgql.UpdateMeetingStateInput{Meeting_id: meetingID, Action: action})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				fmt.Fprintf(os.Stdout, "state updated: meeting %s action=%s success=%v\n", meetingID, resp.UpdateMeetingState.Action, resp.UpdateMeetingState.Success)
				return nil
			}
			return output.Render(os.Stdout, resp.UpdateMeetingState, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "SUCCESS", Path: "success"},
					{Header: "ACTION", Path: "action"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
