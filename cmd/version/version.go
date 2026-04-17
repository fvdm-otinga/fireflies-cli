// Package version implements the `fireflies version` command.
package version

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

// NewVersionCmd returns the `version` command wired with build-time values.
func NewVersionCmd(version, commit, date string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version, commit, and build date",
		RunE: func(cmd *cobra.Command, _ []string) error {
			jsonFlag, _ := cmd.Root().PersistentFlags().GetBool("json")

			info := versionInfo{
				Version: version,
				Commit:  commit,
				Date:    date,
			}

			if jsonFlag {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "fireflies version %s\ncommit:  %s\nbuilt:   %s\n",
				info.Version, info.Commit, info.Date)
			return nil
		},
	}
}
