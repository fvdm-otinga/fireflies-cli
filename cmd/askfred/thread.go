package askfred

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

// newThreadCmd returns `fireflies askfred thread <id>`.
// GraphQL: AskFredThread (query `askfred_thread`)
func newThreadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "thread <id>",
		Short: "Get an AskFred thread by ID with all messages (GraphQL: askfred_thread)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query AskFredThread($id: String!) {\n  askfred_thread(id: $id) { id title created_at messages { id query answer status created_at } }\n}\n")
				fmt.Fprintf(os.Stdout, `{"id": %q}`+"\n", id)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.AskFredThread(context.Background(), c, id)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			if resp.Askfred_thread == nil {
				return ferr.NotFound(fmt.Sprintf("askfred thread %q not found", id))
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.Askfred_thread, output.RenderOpts{
				Format: f,
				Cols:   columns.AskFredThread(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
