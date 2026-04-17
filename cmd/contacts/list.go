package contacts

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

// newListCmd returns `fireflies contacts list`.
// GraphQL: Contacts (query `contacts`)
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List contacts (GraphQL: contacts)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Contacts {\n  contacts { email name last_meeting_date }\n}\n")
				_, _ = os.Stdout.WriteString("{}\n")
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.Contacts(context.Background(), c)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			contacts := resp.Contacts
			if sh.Skip > 0 && sh.Skip < len(contacts) {
				contacts = contacts[sh.Skip:]
			} else if sh.Skip >= len(contacts) && sh.Skip > 0 {
				contacts = nil
			}
			if sh.Limit > 0 && sh.Limit < len(contacts) {
				contacts = contacts[:sh.Limit]
			}

			cur := pagination.NewCursor(sh.Skip, sh.Limit, len(contacts))
			env := output.Envelope(contacts, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
				Format: f,
				Cols:   columns.Contact(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
