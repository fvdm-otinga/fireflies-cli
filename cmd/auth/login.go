package auth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
)

func newLoginCmd() *cobra.Command {
	var profile string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Save a Fireflies API key to the config file",
		Long: `Prompt for a Fireflies API key (input is not echoed), verify it via
the Whoami query, and write it to ~/.config/fireflies/config.toml with
file mode 0600 under the given profile (default: "default").`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Resolve profile name: flag > root --profile > "default".
			if profile == "" {
				profile, _ = cmd.Root().PersistentFlags().GetString("profile")
			}
			if profile == "" {
				profile = "default"
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Enter Fireflies API key for profile %q: ", profile)
			apiKey, err := readSecretLine()
			if err != nil {
				return ferr.General("failed to read API key: " + err.Error())
			}
			apiKey = strings.TrimSpace(apiKey)
			if apiKey == "" {
				return ferr.Usage("API key must not be empty")
			}

			// Verify the key by calling Whoami.
			p := config.Profile{APIKey: apiKey}
			c := client.New(p)
			resp, err := ffgql.Whoami(context.Background(), c)
			if err != nil {
				var cli *ferr.CLIError
				if errors.As(err, &cli) && cli.Exit == ferr.ExitAuthError {
					return ferr.Auth("invalid API key — authentication failed")
				}
				return ferr.Auth("could not verify API key: " + err.Error())
			}
			if resp.User == nil {
				return ferr.Auth("invalid API key — no user returned")
			}

			// Persist.
			loader := config.New()
			if err := loader.SetProfile(profile, p, true); err != nil {
				return ferr.General("write config: " + err.Error())
			}

			email := ""
			if resp.User.Email != nil {
				email = *resp.User.Email
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "\nLogged in as %s (profile: %q)\n", email, profile)
			return nil
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "", "Profile name to save the key under (default: \"default\")")
	return cmd
}

// readSecretLine reads a line from stdin without terminal echo when possible.
// On a TTY, it uses golang.org/x/term.ReadPassword to suppress echo so the
// API key never appears on-screen or in scrollback. For piped input (scripts
// / CI), it falls back to a single-line buffered scan — stdin is not a
// terminal so there's no echo to suppress.
func readSecretLine() (string, error) {
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		if err != nil {
			return "", err
		}
		// ReadPassword swallows the Enter keystroke's echo; emit a newline
		// so subsequent output doesn't glue onto the prompt line.
		fmt.Fprintln(os.Stderr)
		return string(b), nil
	}
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", nil
}
