// Package schema compiles the published MCP JSON Schema files
// (mcp-schema/*.schema.json) and validates arbitrary Go values against any
// named $defs entry within them — the reusable round-trip validator SDK
// smoke tests and the mock referee both call.
package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"sync"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/Axemere-LLC/gismo-contracts"
)

// Registry compiles and caches a jsonschema.Schema per (file, $defs name)
// pair, on demand.
type Registry struct {
	mu    sync.Mutex
	docs  map[string]map[string]any
	cache map[string]*jsonschema.Schema
}

// NewRegistry loads every mcp-schema/*.schema.json file embedded in the
// contracts package.
func NewRegistry() (*Registry, error) {
	entries, err := fs.ReadDir(contracts.MCPSchema, "mcp-schema")
	if err != nil {
		return nil, fmt.Errorf("read mcp-schema dir: %w", err)
	}

	docs := make(map[string]map[string]any, len(entries))
	for _, entry := range entries {
		raw, err := fs.ReadFile(contracts.MCPSchema, "mcp-schema/"+entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		var doc map[string]any
		if err := json.Unmarshal(raw, &doc); err != nil {
			return nil, fmt.Errorf("parse %s: %w", entry.Name(), err)
		}
		docs[entry.Name()] = doc
	}

	return &Registry{docs: docs, cache: make(map[string]*jsonschema.Schema)}, nil
}

// Validate marshals v to JSON and checks it against the named $defs entry
// (e.g. "StateView", "SubmitOrdersResponse") in the given schema file (e.g.
// "getState.schema.json").
func (r *Registry) Validate(file, defName string, v any) error {
	compiled, err := r.compile(file, defName)
	if err != nil {
		return err
	}

	raw, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal value for %s#/$defs/%s: %w", file, defName, err)
	}
	var decoded any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return fmt.Errorf("decode value for %s#/$defs/%s: %w", file, defName, err)
	}

	return compiled.Validate(decoded)
}

// compile builds (or returns the cached) schema for one $defs entry by
// cloning the parsed document and pointing its top-level $ref at that entry
// — the same technique contracts/scripts/validate.sh uses via jq, since
// ajv/jsonschema tooling validates against a document's top-level $ref, not
// an arbitrary JSON-pointer fragment passed alongside the file.
func (r *Registry) compile(file, defName string) (*jsonschema.Schema, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := file + "#/$defs/" + defName
	if cached, ok := r.cache[key]; ok {
		return cached, nil
	}

	doc, ok := r.docs[file]
	if !ok {
		return nil, fmt.Errorf("unknown mcp schema file %q", file)
	}
	defs, ok := doc["$defs"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s has no $defs object", file)
	}
	if _, ok := defs[defName]; !ok {
		return nil, fmt.Errorf("%s has no $defs entry %q", file, defName)
	}

	rewritten := make(map[string]any, len(doc))
	for k, v := range doc {
		rewritten[k] = v
	}
	rewritten["$ref"] = "#/$defs/" + defName

	raw, err := json.Marshal(rewritten)
	if err != nil {
		return nil, fmt.Errorf("marshal %s: %w", key, err)
	}

	// The resource URL itself must not contain a "#" fragment (the compiler
	// panics on newResource otherwise) — encode the $defs pointer into the
	// path instead of a fragment.
	url := "mem://" + file + "/$defs/" + defName
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(url, bytes.NewReader(raw)); err != nil {
		return nil, fmt.Errorf("add resource %s: %w", key, err)
	}
	compiled, err := compiler.Compile(url)
	if err != nil {
		return nil, fmt.Errorf("compile %s: %w", key, err)
	}

	r.cache[key] = compiled
	return compiled, nil
}
