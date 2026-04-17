// Package flags binds the shared CLI flag set to every command.
package flags

import "github.com/spf13/cobra"

// Shared holds the parsed values for flags bound to every command.
type Shared struct {
	Profile    string
	JSON       bool
	JQ         string
	Output     string
	Fields     string
	Limit      int
	Skip       int
	Transcript string
	Since      string
	Until      string
	Yes        bool
	DryRun     bool
}

// Bind attaches the shared flags to cmd. The returned pointer is populated
// when Cobra parses the args.
func Bind(cmd *cobra.Command) *Shared {
	s := &Shared{}
	f := cmd.PersistentFlags()
	f.StringVar(&s.Profile, "profile", "", "Config profile to use")
	f.BoolVar(&s.JSON, "json", false, "Shortcut for --output json")
	f.StringVar(&s.JQ, "jq", "", "Post-process output via a gojq expression")
	f.StringVar(&s.Output, "output", "", "Output format: table|json|ndjson|yaml|tsv|plaintext")
	f.StringVar(&s.Fields, "fields", "", "Comma-separated top-level fields to keep (client-side projection)")
	f.IntVar(&s.Limit, "limit", 0, "Page size (0 = API default, max 50 for transcripts)")
	f.IntVar(&s.Skip, "skip", 0, "Offset pagination cursor")
	f.StringVar(&s.Transcript, "transcript", "", "Transcript depth: none|preview|full")
	f.StringVar(&s.Since, "since", "", "Lower bound (RFC3339 or relative like 7d)")
	f.StringVar(&s.Until, "until", "", "Upper bound (RFC3339)")
	f.BoolVar(&s.Yes, "yes", false, "Bypass confirmation prompts for destructive operations")
	f.BoolVar(&s.DryRun, "dry-run", false, "Print the GraphQL operation without executing")
	return s
}
