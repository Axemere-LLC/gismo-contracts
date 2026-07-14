// Package contracts embeds the versioned OpenAPI contract and MCP JSON
// Schema published from gismo-platform (the source of truth), so the
// conformance harness and downstream Go consumers (e.g. gismo-sdk-go's
// smoke tests) can load them without assuming a filesystem layout relative
// to their own working directory.
package contracts

import "embed"

//go:embed openapi/openapi.yaml
var OpenAPI embed.FS

//go:embed mcp-schema/*.schema.json
var MCPSchema embed.FS

//go:embed examples/*.json
var Examples embed.FS
