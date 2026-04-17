package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// RuleExecutionMeetingGroup returns the standard column set for rule execution meeting groups.
func RuleExecutionMeetingGroup() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "MEETING_ID", Path: "meeting_id"},
		{Header: "MEETING_TITLE", Path: "meeting.title"},
		{Header: "ORGANIZER_EMAIL", Path: "meeting.organizer_email"},
	}
}
