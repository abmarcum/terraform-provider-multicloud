package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNewClientManagerAWSConfig(t *testing.T) {
	model := ProviderModel{
		AWS: &AWSConfigModel{
			Region:    types.StringValue("us-west-2"),
			AccessKey: types.StringValue("AKIAIOSFODNN7EXAMPLE"),
			SecretKey: types.StringValue("wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"),
		},
	}

	ctx := context.Background()
	cm, err := NewClientManager(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error creating ClientManager: %v", err)
	}

	cfg, err := cm.GetAWSConfig(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting AWS config: %v", err)
	}
	if cfg.Region != "us-west-2" {
		t.Errorf("expected AWS region 'us-west-2', got '%s'", cfg.Region)
	}
}

func TestNewClientManagerGCPConfig(t *testing.T) {
	model := ProviderModel{
		GCP: &GCPConfigModel{
			Project: types.StringValue("test-gcp-project"),
			Region:  types.StringValue("us-central1"),
		},
	}

	ctx := context.Background()
	cm, err := NewClientManager(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error creating ClientManager: %v", err)
	}

	cfg, err := cm.GetGCPConfig(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting GCP config: %v", err)
	}
	if cfg.Project != "test-gcp-project" {
		t.Errorf("expected GCP project 'test-gcp-project', got '%s'", cfg.Project)
	}
}

func TestNewClientManagerAzureConfig(t *testing.T) {
	model := ProviderModel{
		Azure: &AzureConfigModel{
			SubscriptionID: types.StringValue("12345678-1234-1234-1234-123456789012"),
			ResourceGroup:  types.StringValue("test-resource-group"),
		},
	}

	ctx := context.Background()
	cm, err := NewClientManager(ctx, model)
	if err != nil {
		t.Fatalf("unexpected error creating ClientManager: %v", err)
	}

	cfg, err := cm.GetAzureConfig(ctx)
	if err != nil {
		t.Fatalf("unexpected error getting Azure config: %v", err)
	}
	if cfg.SubscriptionID != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("expected Azure subscription ID '12345678-1234-1234-1234-123456789012', got '%s'", cfg.SubscriptionID)
	}
}
