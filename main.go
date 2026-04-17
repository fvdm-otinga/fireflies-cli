// Command fireflies is the Fireflies.ai CLI.
package main

import (
	"github.com/fvdm-otinga/fireflies-cli/cmd"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

// Version is injected at build time via -ldflags "-X main.Version=..."
var Version = "dev"

func main() {
	root := cmd.NewRootCmd(Version)
	if err := root.Execute(); err != nil {
		ferr.Exit(err)
	}
}
