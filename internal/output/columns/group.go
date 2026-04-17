package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// UserGroup returns the standard column set for UserGroup responses.
func UserGroup() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "NAME", Path: "name"},
		{Header: "HANDLE", Path: "handle"},
	}
}
