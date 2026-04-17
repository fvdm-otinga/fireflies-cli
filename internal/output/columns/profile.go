package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// Profile returns the standard column set for config profile list output.
func Profile() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "PROFILE", Path: "profile"},
		{Header: "API_KEY", Path: "api_key"},
		{Header: "ENDPOINT", Path: "endpoint"},
		{Header: "ACTIVE", Path: "active"},
	}
}
