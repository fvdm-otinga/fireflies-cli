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

// newShareCmd returns `fireflies meetings share <id> --email <e>...`.
// GraphQL: ShareMeeting (rate-limited: 10/hr via client bucket).
func newShareCmd() *cobra.Command {
	var emails []string

	cmd := &cobra.Command{
		Use:   "share <id>",
		Short: "Share a meeting with one or more email addresses (GraphQL: shareMeeting, rate: 10/hr)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if len(emails) == 0 {
				return ferr.Usage("at least one --email is required")
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{
						"meeting_id": id,
						"emails":     emails,
					},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation ShareMeeting($input: ShareMeetingInput!) {\n  shareMeeting(input: $input) { success message }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.ShareMeeting(context.Background(), c, &ffgql.ShareMeetingInput{
				Meeting_id: id,
				Emails:     emails,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				msg := ""
				if resp.ShareMeeting.Message != nil {
					msg = *resp.ShareMeeting.Message
				}
				fmt.Fprintf(os.Stdout, "meeting %s shared (success=%v) %s\n", id, resp.ShareMeeting.Success, msg)
				return nil
			}
			return output.Render(os.Stdout, resp.ShareMeeting, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "SUCCESS", Path: "success"},
					{Header: "MESSAGE", Path: "message"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringArrayVar(&emails, "email", nil, "Email address to share with (repeatable)")

	return cmd
}
