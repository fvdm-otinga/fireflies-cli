package channels

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

// newGetCmd returns `fireflies channels get <id>`.
// GraphQL: Channel (query `channel`)
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a channel by ID (GraphQL: channel)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Channel($id: ID!) {\n  channel(id: $id) { id title created_at updated_at created_by is_private }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, `{"id": %q}`+"\n", id)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.Channel(context.Background(), c, id)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			if resp.Channel == nil {
				return ferr.NotFound(fmt.Sprintf("channel %q not found", id))
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.Channel, output.RenderOpts{
				Format: f,
				Cols:   columns.Channel(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
