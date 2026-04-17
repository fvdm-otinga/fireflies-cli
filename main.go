// Command fireflies is the Fireflies.ai CLI.
package main

import (
	"github.com/fvdm-otinga/fireflies-cli/cmd"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

// Version, Commit, and Date are injected at build time via
// -ldflags "-X main.Version=... -X main.Commit=... -X main.Date=..."
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func main() {
	root := cmd.NewRootCmd(Version, Commit, Date)
	if err := root.Execute(); err != nil {
		ferr.Exit(err)
	}
}
