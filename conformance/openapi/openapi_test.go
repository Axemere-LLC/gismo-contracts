package openapi

import (
	"context"
	"testing"
)

func TestLoad(t *testing.T) {
	doc, err := Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if doc.Info.Version != "1.0.0" {
		t.Errorf("Info.Version = %q, want %q", doc.Info.Version, "1.0.0")
	}

	const wantOperation = "listTeams"
	path := doc.Paths.Find("/teams")
	if path == nil || path.Get == nil || path.Get.OperationID != wantOperation {
		t.Errorf("GET /teams operationId = %v, want %q", path, wantOperation)
	}
}
