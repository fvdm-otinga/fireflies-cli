package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

func newSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a config value for the active profile",
		Long: `Set a config key for the active (or --profile-specified) profile.
Valid keys: api_key, endpoint.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, value := args[0], args[1]
			profile, _ := cmd.Root().PersistentFlags().GetString("profile")
			if profile == "" {
				profile = "default"
			}

			loader := config.New()
			if err := loader.Load(); err != nil {
				return ferr.General("load config: " + err.Error())
			}

			f, err := loader.All()
			if err != nil {
				return ferr.General("load config: " + err.Error())
			}

			prof := f.Profiles[profile] // zero value if not present

			switch key {
			case "api_key":
				prof.APIKey = value
			case "endpoint":
				prof.Endpoint = value
			default:
				return ferr.Usage(fmt.Sprintf("unknown key %q — valid keys: api_key, endpoint", key))
			}

			if err := loader.SetProfile(profile, prof, false); err != nil {
				return ferr.General("write config: " + err.Error())
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Set %s for profile %q\n", key, profile)
			return nil
		},
	}
	return cmd
}
