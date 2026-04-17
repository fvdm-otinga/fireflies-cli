package channels

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
	"github.com/fvdm-otinga/fireflies-cli/internal/pagination"
)

// newListCmd returns `fireflies channels list`.
// GraphQL: Channels (query `channels`)
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all channels (GraphQL: channels)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Channels {\n  channels { id title created_at updated_at created_by is_private }\n}\n")
				_, _ = os.Stdout.WriteString("{}\n")
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.Channels(context.Background(), c)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			chs := resp.Channels
			if sh.Skip > 0 && sh.Skip < len(chs) {
				chs = chs[sh.Skip:]
			} else if sh.Skip >= len(chs) && sh.Skip > 0 {
				chs = nil
			}
			if sh.Limit > 0 && sh.Limit < len(chs) {
				chs = chs[:sh.Limit]
			}

			cur := pagination.NewCursor(sh.Skip, sh.Limit, len(chs))
			env := output.Envelope(chs, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
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
