package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

var validKeys = []string{"api_key", "endpoint"}

func newGetCmd() *cobra.Command {
	var showSecret bool
	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Print one profile config value",
		Long: `Print the value of a config key for the active profile.
Valid keys: api_key, endpoint.
The api_key is masked unless --show-secret is provided.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key := args[0]
			profile, _ := cmd.Root().PersistentFlags().GetString("profile")
			if profile == "" {
				profile = "default"
			}

			loader := config.New()
			f, err := loader.All()
			if err != nil {
				return ferr.General("load config: " + err.Error())
			}

			prof, ok := f.Profiles[profile]
			if !ok {
				return ferr.NotFound(fmt.Sprintf("profile %q not found in config", profile))
			}

			switch key {
			case "api_key":
				val := prof.APIKey
				if val == "" {
					fmt.Fprintln(cmd.OutOrStdout(), "")
					return nil
				}
				if !showSecret {
					val = maskValue(val)
				}
				fmt.Fprintln(cmd.OutOrStdout(), val)
			case "endpoint":
				fmt.Fprintln(cmd.OutOrStdout(), prof.Endpoint)
			default:
				return ferr.Usage(fmt.Sprintf("unknown key %q — valid keys: api_key, endpoint", key))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&showSecret, "show-secret", false, "Print api_key unmasked")
	return cmd
}

func maskValue(v string) string {
	if len(v) <= 8 {
		return "****"
	}
	return v[:3] + "…" + v[len(v)-4:]
}
