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

// newSoundbiteCmd returns `fireflies live soundbite --meeting <id> --start <t> --end <t>`.
// GraphQL mutation: createLiveSoundbite
func newSoundbiteCmd() *cobra.Command {
	var meetingID, prompt string
	var start, end int

	cmd := &cobra.Command{
		Use:   "soundbite",
		Short: "Create a live soundbite for a meeting (GraphQL: createLiveSoundbite)",
		Long: `Creates a live soundbite for an active meeting.

The --prompt flag describes the soundbite content. Use --start and --end
(seconds) to specify a time window; these are folded into the prompt when
--prompt is not set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)

			if meetingID == "" {
				return ferr.Usage("--meeting is required")
			}

			// Build effective prompt from --start/--end if --prompt not given.
			effectivePrompt := prompt
			if effectivePrompt == "" && (start != 0 || end != 0) {
				effectivePrompt = fmt.Sprintf("soundbite from %ds to %ds", start, end)
			}
			if effectivePrompt == "" {
				return ferr.Usage("--prompt or --start/--end is required")
			}

			if sh.DryRun {
				_, _ = fmt.Fprintf(os.Stdout, "mutation CreateLiveSoundbite($input: CreateLiveSoundbiteInput!) {\n  createLiveSoundbite(input: $input) { success }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, `{"input": {"meeting_id": %q, "prompt": %q}}`+"\n", meetingID, effectivePrompt)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			c := client.New(prof)

			resp, err := ffgql.CreateLiveSoundbite(context.Background(), c, &ffgql.CreateLiveSoundbiteInput{
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
			return output.Render(os.Stdout, resp.CreateLiveSoundbite, output.RenderOpts{
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
	cmd.Flags().StringVar(&prompt, "prompt", "", "Soundbite prompt/description")
	cmd.Flags().IntVar(&start, "start", 0, "Start time in seconds (combined with --end if --prompt not set)")
	cmd.Flags().IntVar(&end, "end", 0, "End time in seconds")
	_ = cmd.MarkFlagRequired("meeting")
	flags.Bind(cmd)
	return cmd
}
