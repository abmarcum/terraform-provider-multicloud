package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type JSONSchema struct {
	Schema      string                 `json:"$schema"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Properties  map[string]PropertyDef `json:"properties"`
	Required    []string               `json:"required"`
}

type PropertyDef struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

func main() {
	schemaData := JSONSchema{
		Schema: "http://json-schema.org/draft-07/schema#",
		Title:  "MultiCloudResourceConfig",
		Type:   "object",
		Properties: map[string]PropertyDef{
			"provider_type": {
				Type:        "string",
				Description: "Target Cloud Provider",
				Enum:        []string{"aws", "gcp", "azure"},
			},
			"resource_name": {
				Type:        "string",
				Description: "Unified resource identifier",
			},
			"region": {
				Type:        "string",
				Description: "Target placement region",
			},
		},
		Required: []string{"provider_type", "resource_name"},
	}

	data, _ := json.MarshalIndent(schemaData, "", "  ")
	/* #nosec G306 */
	_ = os.WriteFile("multicloud_schema.json", data, 0600)
	fmt.Println("[gen-jsonschema] Exported multicloud_schema.json successfully.")
}
