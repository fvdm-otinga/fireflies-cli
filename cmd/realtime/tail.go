package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	"github.com/fvdm-otinga/fireflies-cli/internal/realtime"
)

func newTailCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tail <meeting-id>",
		Short: "Stream live transcript events for a meeting (Socket.IO)",
		Long: `Connects to the Fireflies realtime Socket.IO endpoint and streams
transcript events to stdout. Press Ctrl-C to stop.

Default output is NDJSON (one event per line). Use --format plaintext for
Speaker: text output.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			meetingID := args[0]

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			// Determine output format.
			plaintext := sh.Output == "plaintext"

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle SIGINT/SIGTERM gracefully.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-quit
				cancel()
			}()

			endpoint := realtime.DefaultEndpoint
			c := realtime.New(prof.APIKey, endpoint)

			if err := c.Subscribe(ctx, meetingID, func(e realtime.Event) {
				if plaintext {
					// Attempt to extract Speaker and text from common payload shapes.
					var payload map[string]any
					if err := json.Unmarshal(e.Payload, &payload); err == nil {
						speaker, _ := payload["speaker_name"].(string)
						text, _ := payload["text"].(string)
						if speaker == "" {
							speaker = payload["speaker"].(string)
						}
						if speaker != "" && text != "" {
							_, _ = fmt.Fprintf(os.Stdout, "%s: %s\n", speaker, text)
							return
						}
					}
					// Fallback: print raw payload.
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", e.Payload)
					return
				}
				// NDJSON output.
				line, _ := json.Marshal(map[string]any{
					"event":   e.Name,
					"payload": json.RawMessage(e.Payload),
				})
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", line)
			}); err != nil {
				return ferr.General(fmt.Sprintf("realtime: %v", err))
			}
			return nil
		},
	}
	flags.Bind(cmd)
	return cmd
}
