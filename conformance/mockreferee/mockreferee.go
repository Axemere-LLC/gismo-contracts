// Package mockreferee drives an agent's MCP server through the same fixed
// tool sequence gismo-platform's real referee uses
// (game-and-protocol.md#match-protocol-mcp-tools), validating every response
// against the published MCP JSON Schema. It lets an SDK or reference agent
// prove conformance without needing the private referee.
package mockreferee

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Axemere-LLC/gismo-contracts/conformance/schema"
)

// Step is one canned scenario call: an MCP tool name, the request arguments
// to send, and the schema file/$defs entry its response must validate
// against.
type Step struct {
	Tool            string
	Request         any
	SchemaFile      string
	ResponseDefName string
}

// Scenario is the fixed impulse sequence every agent is driven through: one
// getState round-trip, one submitOrders call, and one surrender poll — the
// same three tools game-and-protocol.md's Turn Loop calls each impulse.
func Scenario(matchID string, impulse int) []Step {
	return []Step{
		{
			Tool: "get_state",
			Request: map[string]any{
				"matchId":      matchID,
				"impulse":      impulse,
				"terrain":      []any{},
				"ownTanks":     []any{},
				"visibleTanks": []any{},
				"blockhouses":  []any{},
			},
			SchemaFile:      "getState.schema.json",
			ResponseDefName: "StateView",
		},
		{
			Tool: "submit_orders",
			Request: map[string]any{
				"matchId": matchID,
				"impulse": impulse,
			},
			SchemaFile:      "submitOrders.schema.json",
			ResponseDefName: "SubmitOrdersResponse",
		},
		{
			Tool: "surrender",
			Request: map[string]any{
				"matchId": matchID,
			},
			SchemaFile:      "surrender.schema.json",
			ResponseDefName: "SurrenderResponse",
		},
	}
}

// Run drives session through every step of the scenario in order, failing
// at the first tool-call error or schema violation.
func Run(ctx context.Context, session *mcp.ClientSession, registry *schema.Registry, steps []Step) error {
	for _, step := range steps {
		result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: step.Tool, Arguments: step.Request})
		if err != nil {
			return fmt.Errorf("call %s: %w", step.Tool, err)
		}
		if result.IsError {
			return fmt.Errorf("agent tool %s returned an error result", step.Tool)
		}
		if err := registry.Validate(step.SchemaFile, step.ResponseDefName, result.StructuredContent); err != nil {
			return fmt.Errorf("%s response failed schema validation: %w", step.Tool, err)
		}
	}
	return nil
}
