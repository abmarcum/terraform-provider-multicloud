package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/abmarcum/multi-cloud-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type JSONSchema struct {
	Schema      string                      `json:"$schema"`
	Title       string                      `json:"title"`
	Description string                      `json:"description"`
	Type        string                      `json:"type"`
	Properties  map[string]PropertyDef      `json:"properties"`
	Required    []string                    `json:"required"`
	Definitions map[string]ResourceTypeSpec `json:"$defs"`
}

type ResourceTypeSpec struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Properties  map[string]PropertyDef `json:"properties"`
	Required    []string               `json:"required,omitempty"`
}

type PropertyDef struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

func inferAttributeJSONType(attr schema.Attribute) string {
	switch attr.(type) {
	case schema.BoolAttribute:
		return "boolean"
	case schema.Int64Attribute, schema.Int32Attribute:
		return "integer"
	case schema.Float64Attribute, schema.Float32Attribute, schema.NumberAttribute:
		return "number"
	case schema.ListAttribute, schema.SetAttribute:
		return "array"
	case schema.MapAttribute, schema.SingleNestedAttribute:
		return "object"
	default:
		return "string"
	}
}

func main() {
	ctx := context.Background()
	p := provider.New("1.0.0")()
	factories := p.Resources(ctx)

	defs := make(map[string]ResourceTypeSpec, len(factories))
	for _, factory := range factories {
		res := factory()
		metaResp := &resource.MetadataResponse{}
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "multicloud"}, metaResp)

		schemaResp := &resource.SchemaResponse{}
		res.Schema(ctx, resource.SchemaRequest{}, schemaResp)

		props := make(map[string]PropertyDef, len(schemaResp.Schema.Attributes))
		var reqAttrs []string

		attrNames := make([]string, 0, len(schemaResp.Schema.Attributes))
		for k := range schemaResp.Schema.Attributes {
			attrNames = append(attrNames, k)
		}
		sort.Strings(attrNames)

		for _, name := range attrNames {
			attr := schemaResp.Schema.Attributes[name]
			prop := PropertyDef{
				Type:        inferAttributeJSONType(attr),
				Description: attr.GetMarkdownDescription(),
			}
			if prop.Description == "" {
				prop.Description = fmt.Sprintf("%s attribute for %s", name, metaResp.TypeName)
			}
			if name == "provider_type" || name == "primary_cloud" || name == "failover_cloud" || name == "source_provider" || name == "destination_provider" {
				prop.Enum = []string{"aws", "gcp", "azure"}
			}
			props[name] = prop
			if attr.IsRequired() {
				reqAttrs = append(reqAttrs, name)
			}
		}

		defs[metaResp.TypeName] = ResourceTypeSpec{
			Type:        "object",
			Description: schemaResp.Schema.Description,
			Properties:  props,
			Required:    reqAttrs,
		}
	}

	schemaData := JSONSchema{
		Schema:      "http://json-schema.org/draft-07/schema#",
		Title:       "MultiCloudResourceConfig",
		Description: fmt.Sprintf("Unified Multi-Cloud Terraform Provider JSON Schema covering %d resources across AWS, GCP, and Azure.", len(defs)),
		Type:        "object",
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
		Required:    []string{"provider_type", "resource_name"},
		Definitions: defs,
	}

	data, _ := json.MarshalIndent(schemaData, "", "  ")
	outPath := filepath.Clean("multicloud_schema.json")
	/* #nosec G306 */
	_ = os.WriteFile(outPath, data, 0600)
	fmt.Printf("[gen-jsonschema] Exported multicloud_schema.json (%d resource definitions) successfully.\n", len(defs))
}
