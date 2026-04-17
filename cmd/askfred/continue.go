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

// newContinueCmd returns `fireflies askfred continue <thread-id> <question>`.
// GraphQL: ContinueAskFredThread.
// Go keyword "continue" is avoided by naming the var continueCmd.
func newContinueCmd() *cobra.Command {
	continueCmd := &cobra.Command{
		Use:   "continue <thread-id> [question]",
		Short: "Continue an existing AskFred conversation thread (GraphQL: continueAskFredThread)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			threadID := args[0]

			question := strings.Join(args[1:], " ")
			if question == "" {
				fi, err := os.Stdin.Stat()
				if err == nil && (fi.Mode()&os.ModeCharDevice) == 0 {
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
						"thread_id": threadID,
						"query":     question,
					},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation ContinueAskFredThread($input: ContinueAskFredThreadInput!) {\n  continueAskFredThread(input: $input) { message { id thread_id query answer status created_at } }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.ContinueAskFredThread(context.Background(), c, &ffgql.ContinueAskFredThreadInput{
				Thread_id: threadID,
				Query:     question,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.ContinueAskFredThread.Message, output.RenderOpts{
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

	flags.Bind(continueCmd)
	return continueCmd
}
