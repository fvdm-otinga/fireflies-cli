package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// Channel returns the standard column set for Channel responses.
func Channel() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "TITLE", Path: "title"},
		{Header: "CREATED_AT", Path: "created_at"},
		{Header: "UPDATED_AT", Path: "updated_at"},
		{Header: "CREATED_BY", Path: "created_by"},
		{Header: "IS_PRIVATE", Path: "is_private"},
	}
}
