package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// AppOutput returns the standard column set for AppOutput responses.
func AppOutput() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "APP_ID", Path: "app_id"},
		{Header: "TITLE", Path: "title"},
		{Header: "TRANSCRIPT_ID", Path: "transcript_id"},
		{Header: "CREATED_AT", Path: "created_at"},
		{Header: "USER_ID", Path: "user_id"},
	}
}
