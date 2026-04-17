package users

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

// newListCmd returns `fireflies users list` — list all workspace users.
// GraphQL: UsersList (query `users`)
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workspace users (GraphQL: users)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query UsersList {\n  users {\n    user_id email name is_admin num_transcripts minutes_consumed\n  }\n}\n")
				_, _ = os.Stdout.WriteString("{}\n")
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.UsersList(context.Background(), c)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			users := resp.Users
			// apply limit/skip client-side (API has no pagination for users)
			if sh.Skip > 0 && sh.Skip < len(users) {
				users = users[sh.Skip:]
			} else if sh.Skip >= len(users) && sh.Skip > 0 {
				users = nil
			}
			if sh.Limit > 0 && sh.Limit < len(users) {
				users = users[:sh.Limit]
			}

			cur := pagination.NewCursor(sh.Skip, sh.Limit, len(users))
			env := output.Envelope(users, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
				Format: f,
				Cols:   columns.User(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
