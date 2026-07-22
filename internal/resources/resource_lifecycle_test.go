package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStorageBucketModelValues(t *testing.T) {
	model := StorageBucketModel{
		ID:                types.StringValue("aws/us-west-2/my-bucket"),
		BucketName:        types.StringValue("my-bucket"),
		ProviderType:      types.StringValue("aws"),
		Region:            types.StringValue("us-west-2"),
		VersioningEnabled: types.BoolValue(true),
		EncryptionEnabled: types.BoolValue(true),
	}

	if model.ID.ValueString() != "aws/us-west-2/my-bucket" {
		t.Errorf("unexpected ID value: %s", model.ID.ValueString())
	}
	if model.ProviderType.ValueString() != "aws" {
		t.Errorf("unexpected ProviderType value: %s", model.ProviderType.ValueString())
	}
	if !model.VersioningEnabled.ValueBool() {
		t.Errorf("expected VersioningEnabled to be true")
	}
}

func TestVirtualMachineModelValues(t *testing.T) {
	model := VirtualMachineModel{
		ID:           types.StringValue("gcp/us-central1/app-vm"),
		VMName:       types.StringValue("app-vm"),
		ProviderType: types.StringValue("gcp"),
		Region:       types.StringValue("us-central1"),
		SizeTier:     types.StringValue("medium"),
	}

	if model.VMName.ValueString() != "app-vm" {
		t.Errorf("unexpected VMName value: %s", model.VMName.ValueString())
	}
	if model.SizeTier.ValueString() != "medium" {
		t.Errorf("unexpected SizeTier value: %s", model.SizeTier.ValueString())
	}
}

func TestResourceConfigureNoop(t *testing.T) {
	bucketRes := NewStorageBucketResource().(*StorageBucketResource)
	req := resource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &resource.ConfigureResponse{}

	bucketRes.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("unexpected error on Configure with nil ProviderData")
	}
}
