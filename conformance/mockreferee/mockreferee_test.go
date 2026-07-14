package mockreferee

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Axemere-LLC/gismo-contracts/conformance/schema"
)

// stateView mirrors internal/referee/wire.go's StateView shape closely
// enough for the scenario's get_state round-trip (echo request back
// unchanged), without depending on the private gismo-platform module.
type stateView struct {
	MatchID      string `json:"matchId"`
	Impulse      int    `json:"impulse"`
	OwnTanks     []any  `json:"ownTanks"`
	VisibleTanks []any  `json:"visibleTanks"`
	Blockhouses  []any  `json:"blockhouses"`
}

type submitOrdersRequest struct {
	MatchID string `json:"matchId"`
	Impulse int    `json:"impulse"`
}

type submitOrdersResponse struct {
	Impulse int   `json:"impulse"`
	Orders  []any `json:"orders"`
}

type surrenderRequest struct {
	MatchID string `json:"matchId"`
}

type surrenderResponse struct {
	Surrendered bool `json:"surrendered"`
}

// connectMockAgent starts an in-process MCP server implementing the three
// match tools and returns a connected client session, with no real network
// involved (mirrors gismo-platform's internal/referee mock-agent test
// pattern).
func connectMockAgent(t *testing.T, ctx context.Context, conformant bool) *mcp.ClientSession {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{Name: "mock-agent", Version: "test"}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "get_state"}, func(_ context.Context, _ *mcp.CallToolRequest, in stateView) (*mcp.CallToolResult, stateView, error) {
		return nil, in, nil
	})

	mcp.AddTool(server, &mcp.Tool{Name: "submit_orders"}, func(_ context.Context, _ *mcp.CallToolRequest, in submitOrdersRequest) (*mcp.CallToolResult, submitOrdersResponse, error) {
		if !conformant {
			// Omit the required "orders" field entirely to violate the schema.
			return nil, submitOrdersResponse{Impulse: in.Impulse, Orders: nil}, nil
		}
		return nil, submitOrdersResponse{Impulse: in.Impulse, Orders: []any{}}, nil
	})

	mcp.AddTool(server, &mcp.Tool{Name: "surrender"}, func(_ context.Context, _ *mcp.CallToolRequest, _ surrenderRequest) (*mcp.CallToolResult, surrenderResponse, error) {
		return nil, surrenderResponse{Surrendered: false}, nil
	})

	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("connect mock agent server: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "gismo-contracts-mockreferee-test", Version: "test"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect mock agent client: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	return session
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		conformant bool
		wantErr    bool
	}{
		{name: "conformant agent passes every step", conformant: true},
		{name: "non-conformant agent fails schema validation", conformant: false, wantErr: true},
	}

	registry, err := schema.NewRegistry()
	if err != nil {
		t.Fatalf("schema.NewRegistry: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			session := connectMockAgent(t, ctx, tt.conformant)

			err := Run(ctx, session, registry, Scenario("test-match", 1))
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
