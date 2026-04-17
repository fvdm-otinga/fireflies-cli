package askfred

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// newAskCmd returns `fireflies askfred ask --meeting <id> <question>`.
// GraphQL: CreateAskFredThread.
// Question may be a positional arg or piped on stdin.
func newAskCmd() *cobra.Command {
	var meetingID string

	cmd := &cobra.Command{
		Use:   "ask [question]",
		Short: "Ask a question about a meeting via AskFred (GraphQL: createAskFredThread)",
		Long: `Ask AskFred a question about a meeting transcript.

The question can be provided as a positional argument or piped via stdin:
  fireflies askfred ask --meeting <id> "summarize in one sentence"
  echo "summarize in one sentence" | fireflies askfred ask --meeting <id>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)

			if meetingID == "" {
				return ferr.Usage("--meeting is required")
			}

			// Collect question from positional args or stdin.
			question := strings.Join(args, " ")
			if question == "" {
				fi, err := os.Stdin.Stat()
				if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
					// stdin has data piped.
					b, err := io.ReadAll(bufio.NewReader(os.Stdin))
					if err != nil {
						return ferr.General(fmt.Sprintf("read stdin: %v", err))
					}
					question = strings.TrimSpace(string(b))
				}
			}
			if question == "" {
				return ferr.Usage("a question is required (positional arg or stdin)")
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{
						"query":         question,
						"transcript_id": meetingID,
					},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation CreateAskFredThread($input: CreateAskFredThreadInput!) {\n  createAskFredThread(input: $input) { message { id thread_id query answer status created_at } }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.CreateAskFredThread(context.Background(), c, &ffgql.CreateAskFredThreadInput{
				Query:         question,
				Transcript_id: &meetingID,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.CreateAskFredThread.Message, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "ID", Path: "id"},
					{Header: "THREAD_ID", Path: "thread_id"},
					{Header: "STATUS", Path: "status"},
					{Header: "QUERY", Path: "query"},
					{Header: "ANSWER", Path: "answer"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&meetingID, "meeting", "", "Meeting (transcript) ID to ask about (required)")

	return cmd
}
