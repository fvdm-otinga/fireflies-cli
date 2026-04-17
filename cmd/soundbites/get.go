package soundbites

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

// newGetCmd returns `fireflies soundbites get <id>`.
// GraphQL: Bite (query `bite`)
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a soundbite by ID (GraphQL: bite)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Bite($id: ID!) {\n  bite(id: $id) { id name transcript_id created_at start_time end_time status }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, `{"id": %q}`+"\n", id)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.Bite(context.Background(), c, id)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			if resp.Bite == nil {
				return ferr.NotFound(fmt.Sprintf("soundbite %q not found", id))
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.Bite, output.RenderOpts{
				Format: f,
				Cols:   columns.Bite(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
