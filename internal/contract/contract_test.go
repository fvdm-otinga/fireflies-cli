package contract_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Khan/genqlient/graphql"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	"github.com/fvdm-otinga/fireflies-cli/internal/contract"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// These tests pin the contract interfaces to their concrete implementations.
// If any command refactor breaks the frozen contract, these fail at compile time.

func TestGraphQLClientInterface(t *testing.T) {
	var _ contract.GraphQLClient = (*client.Client)(nil)
}

func TestConfigLoaderInterface(t *testing.T) {
	var _ contract.ConfigLoader = (*config.Loader)(nil)
}

func TestRendererInterface(t *testing.T) {
	var _ contract.Renderer = rendererShim{}
}

type rendererShim struct{ w *bytes.Buffer }

func (r rendererShim) Render(v any, opts output.RenderOpts) error {
	return output.Render(r.w, v, opts)
}

// Sanity-check that the GraphQL client signature matches what genqlient
// expects — i.e. that client.Client satisfies graphql.Client.
func TestGenqlientClientShape(t *testing.T) {
	var _ graphql.Client = (*client.Client)(nil)
}

// The renderer must accept a plain Go struct and produce JSON without error.
func TestRenderJSONRoundtrip(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}
	var buf bytes.Buffer
	if err := output.Render(&buf, payload{Name: "ok"}, output.RenderOpts{Format: output.FormatJSON}); err != nil {
		t.Fatalf("render: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte(`"name":"ok"`)) {
		t.Fatalf("unexpected output: %s", buf.String())
	}
	_ = context.Background()
}
