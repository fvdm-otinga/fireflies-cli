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

// newActionItemCmd returns `fireflies live action-item --meeting <id> --text <t>`.
// GraphQL mutation: createLiveActionItem
func newActionItemCmd() *cobra.Command {
	var meetingID, prompt, text, assignee string

	cmd := &cobra.Command{
		Use:   "action-item",
		Short: "Create a live action item for a meeting (GraphQL: createLiveActionItem)",
		Long: `Creates a live action item for an active meeting.

Use --prompt (or --text as an alias) for the action item description.
--assignee is recorded in the prompt if provided.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)

			if meetingID == "" {
				return ferr.Usage("--meeting is required")
			}

			// Resolve effective prompt: --text is an alias for --prompt.
			effectivePrompt := prompt
			if effectivePrompt == "" {
				effectivePrompt = text
			}
			if assignee != "" {
				if effectivePrompt != "" {
					effectivePrompt = fmt.Sprintf("%s (assignee: %s)", effectivePrompt, assignee)
				} else {
					effectivePrompt = fmt.Sprintf("action item for %s", assignee)
				}
			}
			if effectivePrompt == "" {
				return ferr.Usage("--prompt or --text is required")
			}

			if sh.DryRun {
				fmt.Fprintf(os.Stdout, "mutation CreateLiveActionItem($input: CreateLiveActionItemInput!) {\n  createLiveActionItem(input: $input) { success }\n}\n")
				fmt.Fprintf(os.Stdout, `{"input": {"meeting_id": %q, "prompt": %q}}`+"\n", meetingID, effectivePrompt)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			c := client.New(prof)

			resp, err := ffgql.CreateLiveActionItem(context.Background(), c, &ffgql.CreateLiveActionItemInput{
				Meeting_id: meetingID,
				Prompt:     effectivePrompt,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.CreateLiveActionItem, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "SUCCESS", Path: "success"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	cmd.Flags().StringVar(&meetingID, "meeting", "", "Meeting ID (required)")
	cmd.Flags().StringVar(&prompt, "prompt", "", "Action item prompt/description")
	cmd.Flags().StringVar(&text, "text", "", "Action item text (alias for --prompt)")
	cmd.Flags().StringVar(&assignee, "assignee", "", "Assignee email (appended to prompt)")
	_ = cmd.MarkFlagRequired("meeting")
	flags.Bind(cmd)
	return cmd
}
