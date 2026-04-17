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

// FromCmd reads the shared flag values from a command that has already had
// Bind called (either on itself or an ancestor via PersistentFlags).
func FromCmd(cmd *cobra.Command) *Shared {
	s := &Shared{}
	get := func(name string) string {
		v, _ := cmd.Flags().GetString(name)
		if v == "" {
			v, _ = cmd.InheritedFlags().GetString(name)
		}
		return v
	}
	getBool := func(name string) bool {
		v, err := cmd.Flags().GetBool(name)
		if err != nil {
			v, _ = cmd.InheritedFlags().GetBool(name)
		}
		return v
	}
	getInt := func(name string) int {
		v, err := cmd.Flags().GetInt(name)
		if err != nil {
			v, _ = cmd.InheritedFlags().GetInt(name)
		}
		return v
	}
	s.Profile = get("profile")
	s.JSON = getBool("json")
	s.JQ = get("jq")
	s.Output = get("output")
	s.Fields = get("fields")
	s.Limit = getInt("limit")
	s.Skip = getInt("skip")
	s.Transcript = get("transcript")
	s.Since = get("since")
	s.Until = get("until")
	s.Yes = getBool("yes")
	s.DryRun = getBool("dry-run")
	return s
}
