package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

func newLogoutCmd() *cobra.Command {
	var (
		profile string
		yes     bool
	)
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove a profile from the config file",
		Long:  `Remove the named profile (default: active profile) from ~/.config/fireflies/config.toml. Use --yes to skip the confirmation prompt.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Resolve profile name.
			if profile == "" {
				profile, _ = cmd.Root().PersistentFlags().GetString("profile")
			}
			if profile == "" {
				profile = "default"
			}
			if !yes {
				yes, _ = cmd.Root().PersistentFlags().GetBool("yes")
			}

			if !yes {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Remove profile %q? [y/N] ", profile)
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
					if answer != "y" && answer != "yes" {
						_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Aborted.")
						return nil
					}
				}
			}

			loader := config.New()
			if err := loader.DeleteProfile(profile); err != nil {
				return ferr.General("remove profile: " + err.Error())
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Profile %q removed.\n", profile)
			return nil
		},
	}
	cmd.Flags().StringVar(&profile, "profile", "", "Profile to remove (default: active profile)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Bypass confirmation prompt")
	return cmd
}
