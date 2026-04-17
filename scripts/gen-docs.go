//go:build gendocs

// gen-docs generates Markdown documentation for all fireflies commands
// using Cobra's built-in doc generator. Output is written to docs/reference/.
//
// Usage:
//
//	go run -tags gendocs scripts/gen-docs.go
package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra/doc"

	"github.com/fvdm-otinga/fireflies-cli/cmd"
)

func main() {
	// Resolve docs/reference/ relative to the repository root regardless of
	// where the script is invoked from.
	_, scriptFile, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("cannot determine script path")
	}
	repoRoot := filepath.Dir(filepath.Dir(scriptFile))
	outDir := filepath.Join(repoRoot, "docs", "reference")

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", outDir, err)
	}

	root := cmd.NewRootCmd("dev", "none", "unknown")
	// Disable the completion command's auto-generated header so it doesn't
	// include machine-specific paths.
	root.DisableAutoGenTag = true

	if err := doc.GenMarkdownTree(root, outDir); err != nil {
		log.Fatalf("doc.GenMarkdownTree: %v", err)
	}

	log.Printf("docs written to %s", outDir)
}
