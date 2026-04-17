package meetings

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	"github.com/fvdm-otinga/fireflies-cli/internal/graphql/dynamic"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newGetCmd returns `fireflies meetings get <id>`.
// GraphQL: DynamicTranscript (query `transcript`) via the dynamic builder.
func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a single meeting/transcript by ID (GraphQL: transcript)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			// Handle --fields ?
			if sh.Fields == "?" {
				_, _ = fmt.Fprintln(os.Stdout, "Available transcript fields:")
				for _, f := range dynamic.AllFields() {
					_, _ = fmt.Fprintf(os.Stdout, "  %s\n", f)
				}
				return nil
			}

			tFields := resolveTranscriptFields(sh.Transcript, sh.Fields)

			if sh.DryRun {
				q := dynamic.BuildSingleTranscriptQuery(tFields)
				_, _ = fmt.Fprintln(os.Stdout, q)
				_, _ = fmt.Fprintf(os.Stdout, `{"id": %q}`+"\n", id)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			result, err := dynamic.ExecuteSingleTranscript(context.Background(), c, id, tFields)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}
			if result == nil {
				return ferr.NotFound(fmt.Sprintf("transcript %q not found", id))
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, result, output.RenderOpts{
				Format: f,
				Cols:   columns.Transcript(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}
