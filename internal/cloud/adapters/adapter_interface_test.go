package adapters

import (
	"context"
	"testing"
)

func TestCloudAdapters(t *testing.T) {
	ctx := context.Background()
	req := ResourceRequest{
		ResourceName: "my-resource",
		ResourceType: "storage_bucket",
		ProviderType: "aws",
		Region:       "us-west-2",
	}

	// 1. AWS Adapter
	aws := &AWSAdapter{}
	resp, err := aws.CreateResource(ctx, req)
	if err != nil || resp.Status != "ACTIVE" {
		t.Errorf("expected AWS adapter create to return ACTIVE status")
	}

	// 2. GCP Adapter
	gcp := &GCPAdapter{}
	resp, err = gcp.CreateResource(ctx, req)
	if err != nil || resp.Status != "RUNNING" {
		t.Errorf("expected GCP adapter create to return RUNNING status")
	}

	// 3. Azure Adapter
	azure := &AzureAdapter{}
	resp, err = azure.CreateResource(ctx, req)
	if err != nil || resp.Status != "SUCCEEDED" {
		t.Errorf("expected Azure adapter create to return SUCCEEDED status")
	}
}
