# gismo-contracts

Public repo, created private-first (flips public at a later reveal milestone).

The published, versioned OpenAPI contract for the Control-Plane API and the MCP tool-surface JSON
Schema, plus the conformance harness used to test SDKs and agents without needing the private referee.
`gismo-platform` is the source of truth for the contract; this repo receives it via a generation
pipeline and is what the SDK repos actually build against.

## Layout

- `openapi/openapi.yaml` — versioned OpenAPI 3.0.3 contract for the Control-Plane REST API
- `mcp-schema/*.schema.json` — MCP tool-surface JSON Schema (draft 2020-12) for `get_state`,
  `submit_orders`, `surrender`
- `examples/*.json` — canned example payloads validated against the schemas above
- `contracts.go` — embeds the three directories above (`OpenAPI`, `MCPSchema`, `Examples`) so Go
  consumers can load the contract without assuming a filesystem layout
- `conformance/` — Go conformance harness:
  - `conformance/openapi` — loads and validates the embedded OpenAPI document
  - `conformance/schema` — compiles the embedded MCP JSON Schema and validates arbitrary values
    against any named `$defs` entry (the reusable round-trip validator)
  - `conformance/mockreferee` — drives an agent's MCP server through the same fixed
    `get_state`/`submit_orders`/`surrender` sequence the real referee uses, validating every
    response against the published schema

## Using this repo

Go consumers (e.g. `gismo-sdk-go`'s smoke tests) can depend on this module directly:

```go
import "github.com/Axemere-LLC/gismo-contracts/conformance/schema"

registry, err := schema.NewRegistry()
err = registry.Validate("getState.schema.json", "StateView", someValue)
```

Non-Go SDKs validate generated models against the same `mcp-schema/*.schema.json` files using their
own language's JSON Schema library.

## License

Apache 2.0 — see `LICENSE`.

## Status

Contract published from `gismo-platform` (source of truth) and conformance harness built, per Phase 3
of the roadmap (see `implementation-roadmap.md` in `gismo-platform`). No registry publish and no CI
workflows yet — both deferred (SDK publishing gates on a later reveal milestone; CI/CD lands in
Phase 7).
