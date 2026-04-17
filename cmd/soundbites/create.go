package soundbites

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newCreateCmd returns `fireflies soundbites create --meeting <id> --start <t> --end <t>`.
// GraphQL: CreateBite.
// Time accepts millisecond epoch or mm:ss format.
func newCreateCmd() *cobra.Command {
	var (
		meetingID string
		startStr  string
		endStr    string
		name      string
		mediaType string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a soundbite from a meeting (GraphQL: createBite)",
		Long: `Create a soundbite clip from a meeting transcript.

Time values (--start, --end) accept:
  - Millisecond epoch:   30000
  - mm:ss format:        0:30 or 1:15`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sh := flags.FromCmd(cmd)

			if meetingID == "" {
				return ferr.Usage("--meeting is required")
			}
			if startStr == "" {
				return ferr.Usage("--start is required")
			}
			if endStr == "" {
				return ferr.Usage("--end is required")
			}

			startMs, err := parseTimeArg(startStr)
			if err != nil {
				return ferr.Usage(fmt.Sprintf("invalid --start: %v", err))
			}
			endMs, err := parseTimeArg(endStr)
			if err != nil {
				return ferr.Usage(fmt.Sprintf("invalid --end: %v", err))
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"transcript_Id": meetingID,
					"name":          name,
					"start_time":    startMs,
					"end_time":      endMs,
					"media_type":    mediaType,
				}, "", "  ")
				_, _ = fmt.Fprintf(os.Stdout, "mutation CreateBite($transcript_Id: ID!, $name: String, $start_time: Float!, $end_time: Float!, $media_type: String) {\n  createBite(transcript_Id: $transcript_Id, name: $name, start_time: $start_time, end_time: $end_time, media_type: $media_type) { id name transcript_id start_time end_time status media_type created_at }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			var namePt *string
			if name != "" {
				namePt = &name
			}
			var mediaTypePt *string
			if mediaType != "" {
				mediaTypePt = &mediaType
			}

			resp, err := ffgql.CreateBite(context.Background(), c, meetingID, namePt, startMs, endMs, mediaTypePt)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.CreateBite, output.RenderOpts{
				Format: f,
				Cols:   columns.Bite(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&meetingID, "meeting", "", "Meeting (transcript) ID (required)")
	cmd.Flags().StringVar(&startStr, "start", "", "Start time in ms or mm:ss (required)")
	cmd.Flags().StringVar(&endStr, "end", "", "End time in ms or mm:ss (required)")
	cmd.Flags().StringVar(&name, "name", "", "Name for the soundbite")
	cmd.Flags().StringVar(&mediaType, "media-type", "", "Media type (audio or video)")

	return cmd
}

// parseTimeArg parses a time argument as either a millisecond epoch integer
// or a mm:ss formatted string, returning float64 milliseconds.
func parseTimeArg(s string) (float64, error) {
	// Try direct numeric (ms epoch).
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v, nil
	}
	// Try mm:ss format.
	parts := strings.SplitN(s, ":", 2)
	if len(parts) == 2 {
		mins, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes in %q", s)
		}
		secs, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid seconds in %q", s)
		}
		return (mins*60 + secs) * 1000, nil
	}
	return 0, fmt.Errorf("expected millisecond epoch or mm:ss, got %q", s)
}
