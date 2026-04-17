package soundbites

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

// newListCmd returns `fireflies soundbites list`.
// GraphQL: Bites (query `bites`)
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List soundbites (GraphQL: bites)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Bites($mine: Boolean, $limit: Int, $skip: Int) {\n  bites(mine: $mine, limit: $limit, skip: $skip) { id name transcript_id created_at start_time end_time status }\n}\n")
				_, _ = os.Stdout.WriteString("{}\n")
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			var limit, skip *int
			if sh.Limit > 0 {
				l := sh.Limit
				limit = &l
			}
			if sh.Skip > 0 {
				s := sh.Skip
				skip = &s
			}

			resp, err := ffgql.Bites(context.Background(), c, nil, nil, limit, skip)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			bites := resp.Bites
			cur := pagination.NewCursor(sh.Skip, sh.Limit, len(bites))
			env := output.Envelope(bites, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
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
