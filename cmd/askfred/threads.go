package askfred

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

// newThreadsCmd returns `fireflies askfred threads`.
// GraphQL: AskFredThreads (query `askfred_threads`)
func newThreadsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threads",
		Short: "List AskFred conversation threads (GraphQL: askfred_threads)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query AskFredThreads {\n  askfred_threads { id title created_at transcript_id user_id }\n}\n")
				_, _ = os.Stdout.WriteString("{}\n")
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.AskFredThreads(context.Background(), c, nil)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			threads := resp.Askfred_threads
			if sh.Skip > 0 && sh.Skip < len(threads) {
				threads = threads[sh.Skip:]
			} else if sh.Skip >= len(threads) && sh.Skip > 0 {
				threads = nil
			}
			if sh.Limit > 0 && sh.Limit < len(threads) {
				threads = threads[:sh.Limit]
			}

			cur := pagination.NewCursor(sh.Skip, sh.Limit, len(threads))
			env := output.Envelope(threads, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
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
