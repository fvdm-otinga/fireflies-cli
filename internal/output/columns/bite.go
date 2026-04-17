package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// Bite returns the standard column set for Bite responses.
func Bite() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "ID", Path: "id"},
		{Header: "NAME", Path: "name"},
		{Header: "TRANSCRIPT_ID", Path: "transcript_id"},
		{Header: "CREATED_AT", Path: "created_at"},
		{Header: "START_TIME", Path: "start_time"},
		{Header: "END_TIME", Path: "end_time"},
		{Header: "STATUS", Path: "status"},
		{Header: "MEDIA_TYPE", Path: "media_type"},
	}
}
