package askfred

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/confirm"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newDeleteCmd returns `fireflies askfred delete <thread-id>`.
// GraphQL: DeleteAskFredThread (destructive → --yes).
func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <thread-id>",
		Short: "Delete an AskFred thread and all its messages (GraphQL: deleteAskFredThread) [destructive]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			threadID := args[0]

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"id": threadID,
				}, "", "  ")
				_, _ = fmt.Fprintf(os.Stdout, "mutation DeleteAskFredThread($id: String!) {\n  deleteAskFredThread(id: $id) { id title created_at }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			if err := confirm.Require(sh.Yes, os.Stdin, fmt.Sprintf("Permanently delete AskFred thread %s and all its messages.", threadID)); err != nil {
				return err
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.DeleteAskFredThread(context.Background(), c, threadID)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				_, _ = fmt.Fprintf(os.Stdout, "thread %s deleted\n", resp.DeleteAskFredThread.Id)
				return nil
			}
			return output.Render(os.Stdout, resp.DeleteAskFredThread, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "ID", Path: "id"},
					{Header: "TITLE", Path: "title"},
					{Header: "CREATED_AT", Path: "created_at"},
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
