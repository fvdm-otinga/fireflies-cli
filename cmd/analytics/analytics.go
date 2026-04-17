// Package analytics implements the `fireflies analytics` command.
package analytics

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

// NewAnalyticsCmd returns the `analytics` command.
// GraphQL: Analytics (query `analytics`)
func NewAnalyticsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics",
		Short: "Fetch workspace analytics (GraphQL: analytics)",
		Long:  `Retrieve team and per-user meeting analytics. Use --since and --until to scope the time window.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if sh.DryRun {
				_, _ = os.Stdout.WriteString("query Analytics($start_time: String, $end_time: String) {\n  analytics(start_time: $start_time, end_time: $end_time) { team { meeting { count duration } } users { user_id user_name user_email } }\n}\n")
				fmt.Fprintf(os.Stdout, `{"start_time": %q, "end_time": %q}`+"\n", sh.Since, sh.Until)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			var since, until *string
			if sh.Since != "" {
				since = &sh.Since
			}
			if sh.Until != "" {
				until = &sh.Until
			}

			resp, err := ffgql.Analytics(context.Background(), c, since, until)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}

			if resp.Analytics == nil {
				return ferr.General("analytics response was empty")
			}
			// For table output, render the users slice for readability.
			// For JSON/YAML, render the full analytics object.
			if f == output.FormatTable || f == output.FormatTSV {
				return output.Render(os.Stdout, resp.Analytics.Users, output.RenderOpts{
					Format: f,
					Cols:   columns.UserAnalytics(),
					Fields: sh.Fields,
					JQ:     sh.JQ,
					Pretty: sh.JSON,
				})
			}
			return output.Render(os.Stdout, resp.Analytics, output.RenderOpts{
				Format: f,
				Cols:   columns.UserAnalytics(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
