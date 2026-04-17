package config

import (
	"os"
	"sort"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all profiles and their config values",
		Long:  `List all profiles in the config file. The api_key is masked.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonFlag, _ := cmd.Root().PersistentFlags().GetBool("json")
			outFmt, _ := cmd.Root().PersistentFlags().GetString("output")
			jq, _ := cmd.Root().PersistentFlags().GetString("jq")
			fields, _ := cmd.Root().PersistentFlags().GetString("fields")

			loader := config.New()
			f, err := loader.All()
			if err != nil {
				return ferr.General("load config: " + err.Error())
			}

			type profileRow struct {
				Profile  string `json:"profile"`
				APIKey   string `json:"api_key"`
				Endpoint string `json:"endpoint"`
				Active   bool   `json:"active"`
			}

			names := make([]string, 0, len(f.Profiles))
			for n := range f.Profiles {
				names = append(names, n)
			}
			sort.Strings(names)

			rows := make([]profileRow, 0, len(names))
			for _, n := range names {
				p := f.Profiles[n]
				masked := ""
				if p.APIKey != "" {
					masked = maskValue(p.APIKey)
				}
				rows = append(rows, profileRow{
					Profile:  n,
					APIKey:   masked,
					Endpoint: p.Endpoint,
					Active:   n == f.Active,
				})
			}

			fmt, err := output.ParseFormat(outFmt, jsonFlag)
			if err != nil {
				return ferr.Usage(err.Error())
			}

			return output.Render(os.Stdout, rows, output.RenderOpts{
				Format: fmt,
				Cols:   columns.Profile(),
				Fields: fields,
				JQ:     jq,
				Pretty: jsonFlag,
			})
		},
	}
	return cmd
}
