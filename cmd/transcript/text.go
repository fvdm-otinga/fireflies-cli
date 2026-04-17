package transcript

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	"github.com/fvdm-otinga/fireflies-cli/internal/graphql/dynamic"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newTextCmd returns `fireflies transcript text <id>`.
// Uses the dynamic builder to request only sentences fields.
// GraphQL: DynamicTranscript with sentences selection.
func newTextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "text <id>",
		Short: "Print meeting transcript as plaintext (GraphQL: transcript.sentences)",
		Long: `Fetch and print the full transcript text for a meeting in Speaker: text format.

  --format plaintext  (default) — one line per sentence: "Speaker: text"
  --format json       — raw sentences JSON array
  --since / --until   — filter sentences by start_time (RFC3339 window)`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			// Dynamic builder: always request sentences + a minimal set
			tFields := []dynamic.TranscriptField{
				dynamic.FID,
				dynamic.FTitle,
				dynamic.FSentences,
			}

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

			// Extract sentences from the raw result
			sentencesRaw, ok := result["sentences"]
			if !ok {
				return ferr.General("no sentences in response")
			}

			// Parse the time filter window (ms epoch from RFC3339)
			var sinceMs, untilMs float64
			if sh.Since != "" {
				if t, err := time.Parse(time.RFC3339, sh.Since); err == nil {
					sinceMs = float64(t.UnixMilli())
				}
			}
			if sh.Until != "" {
				if t, err := time.Parse(time.RFC3339, sh.Until); err == nil {
					untilMs = float64(t.UnixMilli())
				}
			}

			// Decode sentences
			sentencesJSON, err := json.Marshal(sentencesRaw)
			if err != nil {
				return ferr.General("failed to marshal sentences: " + err.Error())
			}
			var sentences []map[string]any
			if err := json.Unmarshal(sentencesJSON, &sentences); err != nil {
				return ferr.General("failed to decode sentences: " + err.Error())
			}

			// Apply time window filter
			filtered := filterSentences(sentences, sinceMs, untilMs)

			// Output format
			outFmt := sh.Output
			if outFmt == "" && !sh.JSON {
				outFmt = "plaintext"
			}
			f, err := output.ParseFormat(outFmt, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}

			if f == output.FormatPlaintext {
				return writeSentencesPlaintext(os.Stdout, filtered)
			}

			// For JSON/YAML/etc render through standard renderer
			return output.Render(os.Stdout, filtered, output.RenderOpts{
				Format: f,
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}
	flags.Bind(cmd)
	return cmd
}

func filterSentences(sentences []map[string]any, sinceMs, untilMs float64) []map[string]any {
	if sinceMs == 0 && untilMs == 0 {
		return sentences
	}
	out := make([]map[string]any, 0, len(sentences))
	for _, s := range sentences {
		startRaw := s["start_time"]
		var start float64
		switch v := startRaw.(type) {
		case float64:
			start = v
		case json.Number:
			start, _ = v.Float64()
		}
		if sinceMs > 0 && start < sinceMs {
			continue
		}
		if untilMs > 0 && start > untilMs {
			continue
		}
		out = append(out, s)
	}
	return out
}

func writeSentencesPlaintext(w io.Writer, sentences []map[string]any) error {
	for _, s := range sentences {
		speaker, _ := s["speaker_name"].(string)
		text, _ := s["text"].(string)
		if speaker == "" {
			speaker = "Unknown"
		}
		if _, err := fmt.Fprintf(w, "%s: %s\n", speaker, text); err != nil {
			return err
		}
	}
	return nil
}
