package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCloudResourceDataSourceMetadataAndSchema(t *testing.T) {
	ctx := context.Background()
	ds := NewCloudResourceDataSource().(*CloudResourceDataSource)

	metaReq := datasource.MetadataRequest{ProviderTypeName: "multicloud"}
	metaResp := &datasource.MetadataResponse{}
	ds.Metadata(ctx, metaReq, metaResp)

	if metaResp.TypeName != "multicloud_resource" {
		t.Errorf("expected TypeName 'multicloud_resource', got '%s'", metaResp.TypeName)
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(ctx, schemaReq, schemaResp)

	for _, attr := range []string{"id", "provider_type", "resource_type", "resource_name", "region", "status", "attributes"} {
		if _, ok := schemaResp.Schema.Attributes[attr]; !ok {
			t.Errorf("multicloud_resource data source missing expected attribute '%s'", attr)
		}
	}

	confReq := datasource.ConfigureRequest{ProviderData: nil}
	confResp := &datasource.ConfigureResponse{}
	ds.Configure(ctx, confReq, confResp)
	if confResp.Diagnostics.HasError() {
		t.Errorf("unexpected error configuring CloudResourceDataSource: %v", confResp.Diagnostics)
	}
}

func TestCostEstimateDataSourceMetadataAndSchema(t *testing.T) {
	ctx := context.Background()
	ds := NewCostEstimateDataSource().(*CostEstimateDataSource)

	metaReq := datasource.MetadataRequest{ProviderTypeName: "multicloud"}
	metaResp := &datasource.MetadataResponse{}
	ds.Metadata(ctx, metaReq, metaResp)

	if metaResp.TypeName != "multicloud_cost_estimate" {
		t.Errorf("expected TypeName 'multicloud_cost_estimate', got '%s'", metaResp.TypeName)
	}

	schemaReq := datasource.SchemaRequest{}
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(ctx, schemaReq, schemaResp)

	for _, attr := range []string{"id", "provider_type", "resource_type", "size_tier", "monthly_cost", "suggested_tier", "estimated_saving", "optimization_note"} {
		if _, ok := schemaResp.Schema.Attributes[attr]; !ok {
			t.Errorf("multicloud_cost_estimate data source missing expected attribute '%s'", attr)
		}
	}
}
