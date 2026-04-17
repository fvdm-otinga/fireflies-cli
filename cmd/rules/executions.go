package rules

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
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newExecutionsCmd returns `fireflies rules executions <meeting-id>`.
// GraphQL: RuleExecutionsByMeeting (query `rule_executions_by_meeting`)
func newExecutionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "executions <meeting-id>",
		Short: "Get rule executions for a meeting (GraphQL: rule_executions_by_meeting)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			meetingID := args[0]

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query RuleExecutionsByMeeting($meeting_id: String, $limit: Int) {\n  rule_executions_by_meeting(filters: { meeting_id: $meeting_id }, limit: $limit) { has_more meetings { meeting_id executions { extension_id extension_title } } }\n}\n")
				fmt.Fprintf(os.Stdout, `{"meeting_id": %q}`+"\n", meetingID)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			var limit *int
			if sh.Limit > 0 {
				l := sh.Limit
				limit = &l
			}

			resp, err := ffgql.RuleExecutionsByMeeting(context.Background(), c, limit, nil, nil, &meetingID)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			var meetings any
			if resp.Rule_executions_by_meeting != nil {
				meetings = resp.Rule_executions_by_meeting.Meetings
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, meetings, output.RenderOpts{
				Format: f,
				Cols:   columns.RuleExecutionMeetingGroup(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
