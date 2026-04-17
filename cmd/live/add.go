package live

import (
	"context"
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

// newAddCmd returns `fireflies live add --meeting <id> --name <n> --email <e>`.
// GraphQL mutation: addToLiveMeeting (rate bucket 3/20min enforced in client).
func newAddCmd() *cobra.Command {
	var meetingLink, name, email string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a participant to a live meeting (GraphQL: addToLiveMeeting, rate: 3/20min)",
		Long: `Adds an attendee to a live/active Fireflies meeting.

The --meeting flag accepts the meeting link (URL) that Fireflies uses to
identify the live session, not the transcript ID.

Rate limit: 3 requests per 20 minutes (enforced client-side).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)

			if meetingLink == "" {
				return ferr.Usage("--meeting is required")
			}

			if sh.DryRun {
				fmt.Fprintf(os.Stdout, "mutation AddToLiveMeeting($meeting_link: String!, ...) {\n  addToLiveMeeting(meeting_link: $meeting_link, ...) { message success }\n}\n")
				fmt.Fprintf(os.Stdout, `{"meeting_link": %q, "name": %q, "email": %q}`+"\n", meetingLink, name, email)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			c := client.New(prof)

			// Build attendees list if name/email provided.
			var attendees []*ffgql.AttendeeInput
			if name != "" || email != "" {
				a := &ffgql.AttendeeInput{}
				if name != "" {
					a.DisplayName = &name
				}
				if email != "" {
					a.Email = &email
				}
				attendees = append(attendees, a)
			}

			resp, err := ffgql.AddToLiveMeeting(context.Background(), c, meetingLink, nil, nil, nil, nil, attendees)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.AddToLiveMeeting, output.RenderOpts{
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

	cmd.Flags().StringVar(&meetingLink, "meeting", "", "Meeting link (URL) of the live meeting (required)")
	cmd.Flags().StringVar(&name, "name", "", "Display name of the attendee to add")
	cmd.Flags().StringVar(&email, "email", "", "Email of the attendee to add")
	_ = cmd.MarkFlagRequired("meeting")
	flags.Bind(cmd)
	return cmd
}
