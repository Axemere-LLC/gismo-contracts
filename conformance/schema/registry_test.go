package schema

import "testing"

func TestRegistry_Validate(t *testing.T) {
	registry, err := NewRegistry()
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}

	tests := []struct {
		name    string
		file    string
		defName string
		value   any
		wantErr bool
	}{
		{
			name:    "getState StateView valid",
			file:    "getState.schema.json",
			defName: "StateView",
			value: map[string]any{
				"matchId":      "m1",
				"impulse":      3,
				"terrain":      []any{},
				"ownTanks":     []any{},
				"visibleTanks": []any{},
				"blockhouses":  []any{},
			},
		},
		{
			name:    "getState StateView missing required field",
			file:    "getState.schema.json",
			defName: "StateView",
			value: map[string]any{
				"matchId": "m1",
				"impulse": 3,
			},
			wantErr: true,
		},
		{
			name:    "getState StateView rejects unknown field",
			file:    "getState.schema.json",
			defName: "StateView",
			value: map[string]any{
				"matchId":      "m1",
				"impulse":      3,
				"terrain":      []any{},
				"ownTanks":     []any{},
				"visibleTanks": []any{},
				"blockhouses":  []any{},
				"extra":        "nope",
			},
			wantErr: true,
		},
		{
			name:    "submitOrders SubmitOrdersResponse valid",
			file:    "submitOrders.schema.json",
			defName: "SubmitOrdersResponse",
			value: map[string]any{
				"impulse": 3,
				"orders": []any{
					map[string]any{
						"tankId": 1, "speed": 2, "heading": 90,
						"turretHold": false, "turretHeading": 100,
						"fire": true, "targetX": 60, "targetY": 44,
					},
				},
			},
		},
		{
			name:    "submitOrders SubmitOrdersResponse wrong type",
			file:    "submitOrders.schema.json",
			defName: "SubmitOrdersResponse",
			value: map[string]any{
				"impulse": "not-an-integer",
				"orders":  []any{},
			},
			wantErr: true,
		},
		{
			name:    "surrender SurrenderResponse valid",
			file:    "surrender.schema.json",
			defName: "SurrenderResponse",
			value:   map[string]any{"surrendered": false},
			wantErr: false,
		},
		{
			name:    "surrender SurrenderResponse missing field",
			file:    "surrender.schema.json",
			defName: "SurrenderResponse",
			value:   map[string]any{},
			wantErr: true,
		},
		{
			name:    "unknown schema file",
			file:    "doesNotExist.schema.json",
			defName: "Whatever",
			value:   map[string]any{},
			wantErr: true,
		},
		{
			name:    "unknown defs entry",
			file:    "surrender.schema.json",
			defName: "DoesNotExist",
			value:   map[string]any{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := registry.Validate(tt.file, tt.defName, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate(%s, %s, %v) error = %v, wantErr %v", tt.file, tt.defName, tt.value, err, tt.wantErr)
			}
		})
	}
}
