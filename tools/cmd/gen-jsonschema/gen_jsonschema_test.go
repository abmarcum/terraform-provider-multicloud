package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJSONSchemaGen(t *testing.T) {
	main()
	defer func() { _ = os.Remove("multicloud_schema.json") }()

	raw, err := os.ReadFile("multicloud_schema.json")
	if err != nil {
		t.Fatalf("expected multicloud_schema.json to be created, got error: %v", err)
	}

	var parsed JSONSchema
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Fatalf("failed to unmarshal multicloud_schema.json: %v", err)
	}
	if len(parsed.Definitions) != 70 {
		t.Errorf("expected 70 resource definitions in JSON Schema, got %d", len(parsed.Definitions))
	}
}
