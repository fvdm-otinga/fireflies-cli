package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// ActiveMeeting returns the standard column set for ActiveMeeting responses.
func ActiveMeeting() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "TITLE", Path: "title"},
		{Header: "START_TIME", Path: "start_time"},
		{Header: "END_TIME", Path: "end_time"},
		{Header: "ORGANIZER_EMAIL", Path: "organizer_email"},
		{Header: "STATE", Path: "state"},
		{Header: "PRIVACY", Path: "privacy"},
		{Header: "MEETING_LINK", Path: "meeting_link"},
	}
}
