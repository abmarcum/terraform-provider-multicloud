package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestStorageBucketResourceMetadata(t *testing.T) {
	r := NewStorageBucketResource()
	req := resource.MetadataRequest{
		ProviderTypeName: "multicloud",
	}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "multicloud_storage_bucket" {
		t.Errorf("expected TypeName 'multicloud_storage_bucket', got '%s'", resp.TypeName)
	}
}

func TestStorageBucketResourceSchema(t *testing.T) {
	r := NewStorageBucketResource()
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if _, ok := resp.Schema.Attributes["bucket_name"]; !ok {
		t.Errorf("expected 'bucket_name' attribute in schema")
	}
	if _, ok := resp.Schema.Attributes["provider_type"]; !ok {
		t.Errorf("expected 'provider_type' attribute in schema")
	}
}
