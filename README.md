# gismo-contracts

Public repo, created private-first (flips public at a later reveal milestone).

The published, versioned OpenAPI contract for the Control-Plane API and the MCP tool-surface JSON
Schema, plus the conformance harness used to test SDKs and agents without needing the private referee.
`gismo-platform` is the source of truth for the contract; this repo receives it via a generation
pipeline and is what the SDK repos actually build against.

## Layout

- `openapi/` — versioned OpenAPI contract (added in a later phase)
- `mcp-schema/` — MCP tool-surface JSON Schema (added in a later phase)
- `conformance/` — canned-scenario mock referee for testing SDKs/agents in isolation

## License

Apache 2.0 — see `LICENSE`.

## Status

Scaffold only. Contract and harness content lands in Phase 3 of the roadmap (see
`implementation-roadmap.md` in `gismo-platform`).
