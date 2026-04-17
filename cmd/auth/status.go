package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
)

type statusInfo struct {
	Profile    string `json:"profile"`
	ConfigFile string `json:"config_file"`
	EnvKeySet  bool   `json:"env_key_set"`
	EnvKeyMask string `json:"env_key_masked,omitempty"`
	Email      string `json:"email,omitempty"`
	Name       string `json:"name,omitempty"`
}

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Print the current authentication status",
		Long: `Print the active profile name, config file path, whether FIREFLIES_API_KEY
is set (masked), and the authenticated user's email (via Whoami).`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profile, _ := cmd.Root().PersistentFlags().GetString("profile")
			jsonFlag, _ := cmd.Root().PersistentFlags().GetBool("json")

			if profile == "" {
				profile = "default"
			}

			loader := config.New()
			configPath, _ := loader.Path()

			// Resolve the env key (masked).
			envKey := os.Getenv(config.EnvAPIKey)
			envKeySet := envKey != ""
			envKeyMask := ""
			if envKeySet {
				envKeyMask = maskKey(envKey)
			}

			// Load the profile (may fail if no key at all).
			prof, profileErr := loader.Profile(profile)

			info := statusInfo{
				Profile:    profile,
				ConfigFile: configPath,
				EnvKeySet:  envKeySet,
				EnvKeyMask: envKeyMask,
			}

			// If we have a key (from env or config), try whoami.
			if profileErr == nil && prof.APIKey != "" {
				c := client.New(prof)
				resp, err := ffgql.Whoami(context.Background(), c)
				if err == nil && resp.User != nil {
					if resp.User.Email != nil {
						info.Email = *resp.User.Email
					}
					if resp.User.Name != nil {
						info.Name = *resp.User.Name
					}
				}
			}

			if jsonFlag {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}

			w := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(w, "Profile:     %s\n", info.Profile)
			_, _ = fmt.Fprintf(w, "Config file: %s\n", info.ConfigFile)
			if info.EnvKeySet {
				_, _ = fmt.Fprintf(w, "Env key:     set (%s)\n", info.EnvKeyMask)
			} else {
				_, _ = fmt.Fprintf(w, "Env key:     not set\n")
			}
			if info.Email != "" {
				_, _ = fmt.Fprintf(w, "Logged in:   %s (%s)\n", info.Email, info.Name)
			} else if profileErr != nil {
				_, _ = fmt.Fprintf(w, "Logged in:   no — %s\n", ferr.General(profileErr.Error()).Message)
			} else {
				_, _ = fmt.Fprintf(w, "Logged in:   unknown (API call failed)\n")
			}
			return nil
		},
	}
	return cmd
}

// maskKey returns a masked form like "ffl…abcd".
func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:3] + "…" + key[len(key)-4:]
}
