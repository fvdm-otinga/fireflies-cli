package meetings

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	"github.com/fvdm-otinga/fireflies-cli/internal/graphql/dynamic"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
	"github.com/fvdm-otinga/fireflies-cli/internal/pagination"
)

// newListCmd returns `fireflies meetings list`.
// GraphQL: DynamicTranscripts (query `transcripts`) via the dynamic builder.
func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List meetings/transcripts (GraphQL: transcripts)",
		Long: `List meetings with pagination. The --transcript flag controls payload depth:
  none    — id, title, date, duration only (default)
  preview — adds organizer_email, participants
  full    — all fields including sentences, summary, analytics

Use --fields ? to list available dynamic fields.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			// Handle --fields ?
			if sh.Fields == "?" {
				_, _ = fmt.Fprintln(os.Stdout, "Available transcript fields:")
				for _, f := range dynamic.AllFields() {
					_, _ = fmt.Fprintf(os.Stdout, "  %s\n", f)
				}
				return nil
			}

			// Resolve selection fields from --transcript and --fields
			tFields := resolveTranscriptFields(sh.Transcript, sh.Fields)

			if sh.DryRun {
				vars := buildListVars(sh)
				q := dynamic.BuildTranscriptsListQuery(tFields)
				_, _ = fmt.Fprintln(os.Stdout, q)
				_, _ = fmt.Fprintf(os.Stdout, "%v\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			limit := sh.Limit
			if limit == 0 {
				limit = 10
			}
			if limit > 50 {
				limit = 50
			}

			vars := buildListVars(sh)
			vars["limit"] = limit
			vars["skip"] = sh.Skip

			result, err := dynamic.ExecuteTranscriptsList(context.Background(), c, vars, tFields)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			cur := pagination.NewCursor(sh.Skip, limit, len(result))
			env := output.Envelope(result, cur.Limit, cur.Skip, cur.NextSkip)

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, env, output.RenderOpts{
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

func buildListVars(sh *flags.Shared) map[string]any {
	vars := map[string]any{}
	if sh.Since != "" {
		if t, err := time.Parse(time.RFC3339, sh.Since); err == nil {
			vars["fromDate"] = t.UTC().Format(time.RFC3339)
		}
	}
	if sh.Until != "" {
		if t, err := time.Parse(time.RFC3339, sh.Until); err == nil {
			vars["toDate"] = t.UTC().Format(time.RFC3339)
		}
	}
	return vars
}

// resolveTranscriptFields converts --transcript / --fields to a []TranscriptField.
func resolveTranscriptFields(transcript, fieldsStr string) []dynamic.TranscriptField {
	// Explicit --fields wins over --transcript
	if fieldsStr != "" && fieldsStr != "?" {
		parsed := dynamic.ParseFields(fieldsStr)
		return dynamic.FieldsFromStrings(parsed)
	}
	switch strings.ToLower(transcript) {
	case "full":
		return dynamic.DefaultFull()
	case "preview":
		return dynamic.DefaultPreview()
	default:
		return dynamic.DefaultNone()
	}
}
