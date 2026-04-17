package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// User returns the standard column set for User responses.
func User() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "USER_ID", Path: "user_id"},
		{Header: "EMAIL", Path: "email"},
		{Header: "NAME", Path: "name"},
		{Header: "ADMIN", Path: "is_admin"},
		{Header: "TRANSCRIPTS", Path: "num_transcripts"},
		{Header: "MINUTES", Path: "minutes_consumed"},
	}
}
