package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// Contact returns the standard column set for Contact responses.
func Contact() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "EMAIL", Path: "email"},
		{Header: "NAME", Path: "name"},
		{Header: "LAST_MEETING_DATE", Path: "last_meeting_date"},
		{Header: "PICTURE", Path: "picture"},
	}
}
