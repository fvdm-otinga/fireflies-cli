package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// AskFredThread returns the standard column set for AskFredThread summary responses.
func AskFredThread() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "TITLE", Path: "title"},
		{Header: "CREATED_AT", Path: "created_at"},
		{Header: "TRANSCRIPT_ID", Path: "transcript_id"},
		{Header: "USER_ID", Path: "user_id"},
	}
}
