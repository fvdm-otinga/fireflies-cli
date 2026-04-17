package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// UserAnalytics returns the standard column set for per-user analytics rows.
func UserAnalytics() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "USER_ID", Path: "user_id"},
		{Header: "USER_NAME", Path: "user_name"},
		{Header: "USER_EMAIL", Path: "user_email"},
		{Header: "MEETINGS", Path: "meeting.count"},
		{Header: "DURATION", Path: "meeting.duration"},
		{Header: "TALK_LISTEN_PCT", Path: "conversation.talk_listen_pct"},
		{Header: "WORD_COUNT", Path: "conversation.total_word_count"},
	}
}
