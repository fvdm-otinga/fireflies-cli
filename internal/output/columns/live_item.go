package columns

import "github.com/fvdm-otinga/fireflies-cli/internal/output"

// LiveActionItem returns the standard column set for LiveActionItem responses.
func LiveActionItem() []output.ColumnDef {
	return []output.ColumnDef{
		{Header: "NAME", Path: "name"},
		{Header: "ACTION_ITEM", Path: "action_item"},
	}
}
