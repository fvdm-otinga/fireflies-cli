package meetings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newMoveCmd returns `fireflies meetings move <id> --channel <cid>`.
// GraphQL: UpdateMeetingChannel
func newMoveCmd() *cobra.Command {
	var channelID string

	cmd := &cobra.Command{
		Use:   "move <id>",
		Short: "Move a meeting to a channel (GraphQL: updateMeetingChannel)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if channelID == "" {
				return ferr.Usage("--channel is required")
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{
						"transcript_ids": []string{id},
						"channel_id":     channelID,
					},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation UpdateMeetingChannel($input: UpdateMeetingChannelInput!) {\n  updateMeetingChannel(input: $input) { id title date channels { id title } }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.UpdateMeetingChannel(context.Background(), c, &ffgql.UpdateMeetingChannelInput{
				Transcript_ids: []string{id},
				Channel_id:     channelID,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				fmt.Fprintf(os.Stdout, "meeting %s moved to channel %s\n", id, channelID)
				return nil
			}
			return output.Render(os.Stdout, resp.UpdateMeetingChannel, output.RenderOpts{
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
	cmd.Flags().StringVar(&channelID, "channel", "", "Channel ID to move the meeting into (required)")

	return cmd
}
