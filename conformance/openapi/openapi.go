// Package openapi loads and validates the published Control-Plane OpenAPI
// document, confirming it parses and is internally consistent (every $ref
// resolves, every schema is well-formed) independent of the openapi-generator
// toolchain used in gismo-platform.
package openapi

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/Axemere-LLC/gismo-contracts"
)

// Load parses and validates openapi/openapi.yaml from the embedded contract.
func Load(ctx context.Context) (*openapi3.T, error) {
	raw, err := contracts.OpenAPI.ReadFile("openapi/openapi.yaml")
	if err != nil {
		return nil, fmt.Errorf("read openapi.yaml: %w", err)
	}

	doc, err := openapi3.NewLoader().LoadFromData(raw)
	if err != nil {
		return nil, fmt.Errorf("parse openapi.yaml: %w", err)
	}
	if err := doc.Validate(ctx); err != nil {
		return nil, fmt.Errorf("validate openapi.yaml: %w", err)
	}

	return doc, nil
}
