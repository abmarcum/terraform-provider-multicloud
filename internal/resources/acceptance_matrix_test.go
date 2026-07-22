package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAcceptanceMatrixSchemaAndStateLifecycle(t *testing.T) {
	ctx := context.Background()

	resourcesToAcceptanceTest := []struct {
		name    string
		factory func() resource.Resource
	}{
		{"multicloud_storage_bucket", NewStorageBucketResource},
		{"multicloud_virtual_machine", NewVirtualMachineResource},
		{"multicloud_virtual_network", NewVirtualNetworkResource},
		{"multicloud_db_instance", NewDBInstanceResource},
		{"multicloud_failover_policy", NewFailoverPolicyResource},
		{"multicloud_data_sync", NewDataSyncResource},
	}

	for _, tc := range resourcesToAcceptanceTest {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.factory()

			// 1. Verify Metadata
			metaReq := resource.MetadataRequest{ProviderTypeName: "multicloud"}
			metaResp := &resource.MetadataResponse{}
			r.Metadata(ctx, metaReq, metaResp)

			if metaResp.TypeName != tc.name {
				t.Errorf("expected TypeName '%s', got '%s'", tc.name, metaResp.TypeName)
			}

			// 2. Verify Schema
			schemaReq := resource.SchemaRequest{}
			schemaResp := &resource.SchemaResponse{}
			r.Schema(ctx, schemaReq, schemaResp)

			if schemaResp.Schema.Attributes == nil {
				t.Fatalf("attributes for %s cannot be nil", tc.name)
			}
		})
	}
}
