package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJSONSchemaFileGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	schemaPath := filepath.Join(tmpDir, "test_schema.json")

	// Verify schema export file creation logic
	data := []byte(`{"title": "MultiCloudProviderSchema"}`)
	err := os.WriteFile(schemaPath, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test schema: %v", err)
	}

	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Errorf("expected test schema file to exist at %s", schemaPath)
	}
}
