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

func TestValidatePreApplySecurityOnUpdateAndPolicyControls(t *testing.T) {
	ctx := context.Background()
	t.Setenv("MULTICLOUD_MOCK_MODE", "true")

	// 1. Encryption-at-rest enforcement on UpdateCloudResource for storage/vault resources
	for _, resType := range []string{"block_volume", "shared_filesystem", "backup_vault"} {
		_, err := UpdateCloudResource(ctx, "aws", resType, "my-res", "us-west-2", map[string]interface{}{
			"encryption_enabled": false,
		})
		if err == nil {
			t.Errorf("expected UpdateCloudResource to reject unencrypted %s when encryption_enabled=false", resType)
		}
	}

	// 2. Secret leak detection on UpdateCloudResource for inline code/config/workflow definitions
	for _, attrKey := range []string{"code_content", "config_value", "definition"} {
		_, err := UpdateCloudResource(ctx, "aws", "edge_function", "my-fn", "us-west-2", map[string]interface{}{
			attrKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...",
		})
		if err == nil {
			t.Errorf("expected UpdateCloudResource to reject RSA private key in %s", attrKey)
		}
	}

	// 3. Invalid provider_type and invalid CIDR block rejection
	if _, err := CreateCloudResource(ctx, "digitalocean", "storage_bucket", "my-bucket", "nyc1", nil); err == nil {
		t.Errorf("expected CreateCloudResource to reject unsupported provider_type 'digitalocean'")
	}
	if _, err := CreateCloudResource(ctx, "aws", "virtual_network", "my-vpc", "us-west-2", map[string]interface{}{
		"cidr_block": "999.0.0.0/8",
	}); err == nil {
		t.Errorf("expected CreateCloudResource to reject malformed cidr_block")
	}

	// 4. Case-insensitive IsSensitiveOrInternalKey checks
	if !IsSensitiveOrInternalKey("AWS_SECRET_KEY") || !IsSensitiveOrInternalKey("_CLIENT_MANAGER") {
		t.Errorf("expected IsSensitiveOrInternalKey to match uppercase sensitive and internal keys")
	}
	if IsSensitiveOrInternalKey("environment_tier") {
		t.Errorf("expected IsSensitiveOrInternalKey to allow benign key 'environment_tier'")
	}
}

