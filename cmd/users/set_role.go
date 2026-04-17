package users

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/confirm"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
	"github.com/fvdm-otinga/fireflies-cli/internal/output/columns"
)

// newSetRoleCmd returns `fireflies users set-role <uid> <admin|user>`.
// GraphQL: SetUserRole (destructive → --yes).
func newSetRoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-role <user-id> <admin|user>",
		Short: "Set a user's role (GraphQL: setUserRole) [destructive]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			userID := args[0]
			roleStr := strings.ToLower(args[1])

			role := ffgql.Role(roleStr)
			switch role {
			case ffgql.RoleAdmin, ffgql.RoleUser:
				// valid
			default:
				return ferr.Usage(fmt.Sprintf("invalid role %q: must be admin or user", roleStr))
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"user_id": userID,
					"role":    roleStr,
				}, "", "  ")
				_, _ = fmt.Fprintf(os.Stdout, "mutation SetUserRole($user_id: String!, $role: Role!) {\n  setUserRole(user_id: $user_id, role: $role) { user_id email name is_admin }\n}\n")
				_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			if err := confirm.Require(sh.Yes, os.Stdin, fmt.Sprintf("Set role of user %s to %s.", userID, roleStr)); err != nil {
				return err
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.SetUserRole(context.Background(), c, userID, role)
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			return output.Render(os.Stdout, resp.SetUserRole, output.RenderOpts{
				Format: f,
				Cols:   columns.User(),
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	flags.Bind(cmd)
	return cmd
}
