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

	// 1. Test Mock Mode (MULTICLOUD_MOCK_MODE=true)
	t.Setenv("MULTICLOUD_MOCK_MODE", "true")

	aws := &AWSAdapter{}
	resp, err := aws.CreateResource(ctx, req)
	if err != nil || resp.Status != "ACTIVE" {
		t.Errorf("expected AWS adapter create in mock mode to return ACTIVE status")
	}

	gcp := &GCPAdapter{}
	resp, err = gcp.CreateResource(ctx, req)
	if err != nil || resp.Status != "RUNNING" {
		t.Errorf("expected GCP adapter create in mock mode to return RUNNING status")
	}

	azure := &AzureAdapter{}
	resp, err = azure.CreateResource(ctx, req)
	if err != nil || resp.Status != "SUCCEEDED" {
		t.Errorf("expected Azure adapter create in mock mode to return SUCCEEDED status")
	}

	// 2. Test Live Mode (MULTICLOUD_MOCK_MODE="") - should return live cloud authentication error when uncredentialed
	t.Setenv("MULTICLOUD_MOCK_MODE", "")
	_, err = azure.CreateResource(ctx, req)
	if err == nil {
		t.Errorf("expected Azure adapter create in live mode without credentials to return authentication error")
	}

	// 3. Test Pre-Apply Secret Leak Prevention in CreateCloudResource
	t.Setenv("MULTICLOUD_MOCK_MODE", "true")
	_, err = CreateCloudResource(ctx, "aws", "storage_bucket", "my-bucket", "us-west-2", map[string]interface{}{
		"leaked_key": "aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	})
	if err == nil {
		t.Errorf("expected CreateCloudResource to reject attribute containing leaked AWS secret access key")
	}

	// 4. Test Sensitive Credential Stripping & Region Sanitization
	safeResp, err := CreateCloudResource(ctx, "aws", "storage_bucket", "my-bucket", "us-west-2@evil.com", map[string]interface{}{
		"aws_secret_key": "should-not-echo-in-response",
		"custom_tag":     "allowed",
	})
	if err != nil {
		t.Fatalf("unexpected error in CreateCloudResource: %v", err)
	}
	if _, leaked := safeResp.Attributes["aws_secret_key"]; leaked {
		t.Errorf("expected aws_secret_key to be stripped from ResourceResponse.Attributes")
	}
	if safeResp.Attributes["region"] != "us-west-2evilcom" {
		t.Errorf("expected region authority characters to be sanitized, got %v", safeResp.Attributes["region"])
	}
}

