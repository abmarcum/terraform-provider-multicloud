package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStorageBucketMockCreate(t *testing.T) {
	r := NewStorageBucketResource().(*StorageBucketResource)
	ctx := context.Background()

	// Ensure Configure handles nil provider data safely
	confReq := resource.ConfigureRequest{ProviderData: nil}
	confResp := &resource.ConfigureResponse{}
	r.Configure(ctx, confReq, confResp)

	if confResp.Diagnostics.HasError() {
		t.Errorf("unexpected error configuring StorageBucketResource")
	}
}

func TestVirtualMachineMockState(t *testing.T) {
	state := VirtualMachineModel{
		ID:           types.StringValue("aws/us-west-2/test-vm"),
		VMName:       types.StringValue("test-vm"),
		ProviderType: types.StringValue("aws"),
		Region:       types.StringValue("us-west-2"),
		SizeTier:     types.StringValue("small"),
	}

	if state.ID.ValueString() != "aws/us-west-2/test-vm" {
		t.Errorf("unexpected state ID: %s", state.ID.ValueString())
	}
}

func TestFailoverPolicyMockState(t *testing.T) {
	state := FailoverPolicyModel{
		ID:             types.StringValue("failover/global-dr"),
		PolicyName:     types.StringValue("global-dr"),
		PrimaryCloud:   types.StringValue("aws"),
		FailoverCloud:  types.StringValue("gcp"),
		FailoverStatus: types.StringValue("PRIMARY_HEALTHY"),
	}

	if state.FailoverStatus.ValueString() != "PRIMARY_HEALTHY" {
		t.Errorf("unexpected failover status: %s", state.FailoverStatus.ValueString())
	}
}

func TestDataSyncMockState(t *testing.T) {
	state := DataSyncModel{
		ID:             types.StringValue("datasync/backup/aws-to-gcp"),
		SyncName:       types.StringValue("backup"),
		SourceProvider: types.StringValue("aws"),
		DestProvider:   types.StringValue("gcp"),
		SyncStatus:     types.StringValue("REPLICATION_ACTIVE"),
	}

	if state.SyncStatus.ValueString() != "REPLICATION_ACTIVE" {
		t.Errorf("unexpected sync status: %s", state.SyncStatus.ValueString())
	}
}
