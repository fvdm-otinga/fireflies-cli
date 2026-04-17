package meetings

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/confirm"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newRevokeCmd returns `fireflies meetings revoke <id> --email <e>`.
// GraphQL: RevokeSharedMeetingAccess (destructive → --yes).
func newRevokeCmd() *cobra.Command {
	var email string

	cmd := &cobra.Command{
		Use:   "revoke <id>",
		Short: "Revoke shared meeting access for an email (GraphQL: revokeSharedMeetingAccess) [destructive]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			id := args[0]

			if email == "" {
				return ferr.Usage("--email is required")
			}

			if sh.DryRun {
				vars, _ := json.MarshalIndent(map[string]any{
					"input": map[string]any{
						"meeting_id": id,
						"email":      email,
					},
				}, "", "  ")
				fmt.Fprintf(os.Stdout, "mutation RevokeSharedMeetingAccess($input: RevokeSharedMeetingAccessInput!) {\n  revokeSharedMeetingAccess(input: $input) { success message }\n}\n")
				fmt.Fprintf(os.Stdout, "%s\n", vars)
				return nil
			}

			if err := confirm.Require(sh.Yes, os.Stdin, fmt.Sprintf("Revoke access for %s from meeting %s.", email, id)); err != nil {
				return err
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			resp, err := ffgql.RevokeSharedMeetingAccess(context.Background(), c, &ffgql.RevokeSharedMeetingAccessInput{
				Meeting_id: id,
				Email:      email,
			})
			if err != nil {
				return ferr.FromGraphQLError(err)
			}

			f, err := output.ParseFormat(sh.Output, sh.JSON)
			if err != nil {
				return ferr.Usage(err.Error())
			}
			if f == output.FormatTable && !sh.JSON {
				msg := ""
				if resp.RevokeSharedMeetingAccess.Message != nil {
					msg = *resp.RevokeSharedMeetingAccess.Message
				}
				fmt.Fprintf(os.Stdout, "access revoked for %s from meeting %s (success=%v) %s\n", email, id, resp.RevokeSharedMeetingAccess.Success, msg)
				return nil
			}
			return output.Render(os.Stdout, resp.RevokeSharedMeetingAccess, output.RenderOpts{
				Format: f,
				Cols: []output.ColumnDef{
					{Header: "SUCCESS", Path: "success"},
					{Header: "MESSAGE", Path: "message"},
				},
				Fields: sh.Fields,
				JQ:     sh.JQ,
				Pretty: sh.JSON,
			})
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&email, "email", "", "Email address to revoke access for (required)")

	return cmd
}
