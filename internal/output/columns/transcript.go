package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// Transcript returns the standard column set for Transcript list responses.
func Transcript() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "TITLE", Path: "title"},
		{Header: "DATE", Path: "date"},
		{Header: "DURATION", Path: "duration"},
		{Header: "HOST_EMAIL", Path: "host_email"},
		{Header: "ORGANIZER_EMAIL", Path: "organizer_email"},
		{Header: "MEETING_LINK", Path: "meeting_link"},
		{Header: "PRIVACY", Path: "privacy"},
	}
}
