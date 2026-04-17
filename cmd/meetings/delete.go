package meetings

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

// newDeleteCmd returns `fireflies meetings delete <id>`.
// GraphQL: DeleteTranscript (destructive → --yes, rate: 10/min via client bucket).
func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a meeting/transcript (GraphQL: deleteTranscript, rate: 10/min) [destructive]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"id": id,
				}, "", "  ")
				_, _ = fmt.Fprintf(os.Stdout, "mutation DeleteTranscript($id: String!) {\n  deleteTranscript(id: $id) { id title }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			if err := confirm.Require(sh.Yes, os.Stdin, fmt.Sprintf("Permanently delete meeting %s.", id)); err != nil {
				return err
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.DeleteTranscript(context.Background(), c, id)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				title := ""
				if resp.DeleteTranscript.Title != nil {
					title = *resp.DeleteTranscript.Title
				}
				_, _ = fmt.Fprintf(os.Stdout, "meeting %s deleted (%s)\n", id, title)
				return nil
			}
			return output.Render(os.Stdout, resp.DeleteTranscript, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "ID", Path: "id"},
					{Header: "TITLE", Path: "title"},
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
