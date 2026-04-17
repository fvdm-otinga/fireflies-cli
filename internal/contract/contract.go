// Package contract defines the interfaces that every command in the CLI
// builds against. These interfaces are frozen at the `contract-v1` git tag;
// changes require a documented RFC in docs/rfc/ and owner approval.
//
// See docs/interface-contract.md for the narrative.
package contract

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// GraphQLClient is implemented by internal/client.Client. Every command uses
// this to talk to the Fireflies API.
type GraphQLClient interface {
	MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error
	Endpoint() string
}

// ConfigLoader is implemented by internal/config.Loader.
type ConfigLoader interface {
	Profile(name string) (config.Profile, error)
	Path() (string, error)
}

// Renderer is implemented by internal/output.Render (wrapped in a shim).
type Renderer interface {
	Render(v any, opts output.RenderOpts) error
}
