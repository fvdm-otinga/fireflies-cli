package users

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newWhoamiCmd returns `fireflies users whoami` — the canonical end-to-end
// smoke command. It exercises auth, the GraphQL client, and every output
// format in the renderer.
//
// GraphQL: Whoami (query `user`)
func newWhoamiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Print the API key owner's profile (GraphQL: user)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			root := cmd.Root()
			profile, _ := root.PersistentFlags().GetString("profile")
			jsonFlag, _ := root.PersistentFlags().GetBool("json")
			outFmt, _ := root.PersistentFlags().GetString("output")
			jq, _ := root.PersistentFlags().GetString("jq")
			fields, _ := root.PersistentFlags().GetString("fields")

			prof, err := config.New().Profile(profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.Whoami(context.Background(), c)
			if err != nil {
				return ferr.General(err.Error())
			}
			if resp.User == nil {
				return ferr.NotFound("no user returned from API (check your token)")
			}

			f, err := output.ParseFormat(outFmt, jsonFlag)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.User, output.RenderOpts{
				Format: f,
				Cols:   columns.User(),
				Fields: fields,
				JQ:     jq,
				Pretty: jsonFlag,
			})
		},
	}
}
