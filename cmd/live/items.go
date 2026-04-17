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
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newItemsCmd returns `fireflies live items <meeting-id>`.
// GraphQL: LiveActionItems (query `live_action_items`)
func newItemsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "items <meeting-id>",
		Short: "Get live action items for a meeting (GraphQL: live_action_items)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			meetingID := args[0]

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query LiveActionItems($meeting_id: ID!) {\n  live_action_items(meeting_id: $meeting_id) { name action_item }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, `{"meeting_id": %q}`+"\n", meetingID)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.LiveActionItems(context.Background(), c, meetingID)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.Live_action_items, output.RenderOpts{
				Format: f,
				Cols:   columns.LiveActionItem(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
