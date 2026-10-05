package main

import (
	"os"
	"testing"
)

func TestJSONSchemaGen(t *testing.T) {
	main()
	defer func() { _ = os.Remove("multicloud_schema.json") }()

	if _, err := os.Stat("multicloud_schema.json"); err != nil {
		t.Fatalf("expected multicloud_schema.json to be created, got error: %v", err)
	}
}
