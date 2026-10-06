package provider

import (
	"context"
	"strings"
	"sync"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var _ common.ClientConfigProvider = &ClientManager{}

// ClientManager holds cloud SDK instances with thread-safe lazy initialization
type ClientManager struct {
	model             ProviderModel
	awsConfig         *aws.Config
	gcpConfig         *GCPClientConfig
	azureConfig       *AzureClientConfig
	awsInitOnce       sync.Once
	gcpInitOnce       sync.Once
	azureInitOnce     sync.Once
	defaultAttrsCache sync.Map
}

type GCPClientConfig struct {
	Project     string
	Region      string
	Credentials string
}

type AzureClientConfig struct {
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
	ResourceGroup  string
}

// NewClientManager returns a ClientManager that defers SDK initialization until clients are actually requested
func NewClientManager(ctx context.Context, model ProviderModel) (*ClientManager, error) {
	cm := &ClientManager{
		model: model,
	}
	return cm, nil
}

// IsMockMode returns true if synthetic mock mode was enabled on the provider configuration or environment
func (cm *ClientManager) IsMockMode() bool {
	if !cm.model.MockMode.IsNull() && !cm.model.MockMode.IsUnknown() && cm.model.MockMode.ValueBool() {
		return true
	}
	return common.IsMockMode()
}

func (cm *ClientManager) computeStaticDefaultAttrs(pType string) map[string]interface{} {
	if cached, ok := cm.defaultAttrsCache.Load(pType); ok {
		return cached.(map[string]interface{})
	}

	base := make(map[string]interface{}, 6)
	if !cm.model.DefaultTags.IsNull() && !cm.model.DefaultTags.IsUnknown() {
		elems := cm.model.DefaultTags.Elements()
		if len(elems) > 0 {
			tagMap := make(map[string]string, len(elems))
			for k, v := range elems {
				tagMap[k] = strings.Trim(v.String(), "\"")
			}
			base["default_tags"] = tagMap
		}
	}

	switch pType {
	case "aws":
		if cm.model.AWS != nil {
			if !cm.model.AWS.AccessKey.IsNull() && !cm.model.AWS.AccessKey.IsUnknown() && cm.model.AWS.AccessKey.ValueString() != "" {
				base["aws_access_key"] = cm.model.AWS.AccessKey.ValueString()
			}
			if !cm.model.AWS.SecretKey.IsNull() && !cm.model.AWS.SecretKey.IsUnknown() && cm.model.AWS.SecretKey.ValueString() != "" {
				base["aws_secret_key"] = cm.model.AWS.SecretKey.ValueString()
			}
			if !cm.model.AWS.Profile.IsNull() && !cm.model.AWS.Profile.IsUnknown() && cm.model.AWS.Profile.ValueString() != "" {
				base["aws_profile"] = cm.model.AWS.Profile.ValueString()
			}
			if !cm.model.AWS.Region.IsNull() && !cm.model.AWS.Region.IsUnknown() && cm.model.AWS.Region.ValueString() != "" {
				base["provider_default_region"] = cm.model.AWS.Region.ValueString()
			}
		}
	case "gcp":
		if cm.model.GCP != nil {
			if !cm.model.GCP.Project.IsNull() && !cm.model.GCP.Project.IsUnknown() && cm.model.GCP.Project.ValueString() != "" {
				base["gcp_project"] = cm.model.GCP.Project.ValueString()
			}
			if !cm.model.GCP.Credentials.IsNull() && !cm.model.GCP.Credentials.IsUnknown() && cm.model.GCP.Credentials.ValueString() != "" {
				base["gcp_credentials"] = cm.model.GCP.Credentials.ValueString()
			}
			if !cm.model.GCP.Region.IsNull() && !cm.model.GCP.Region.IsUnknown() && cm.model.GCP.Region.ValueString() != "" {
				base["provider_default_region"] = cm.model.GCP.Region.ValueString()
			}
		}
	case "azure":
		if cm.model.Azure != nil {
			if !cm.model.Azure.SubscriptionID.IsNull() && !cm.model.Azure.SubscriptionID.IsUnknown() && cm.model.Azure.SubscriptionID.ValueString() != "" {
				base["azure_subscription_id"] = cm.model.Azure.SubscriptionID.ValueString()
			}
			if !cm.model.Azure.ResourceGroup.IsNull() && !cm.model.Azure.ResourceGroup.IsUnknown() && cm.model.Azure.ResourceGroup.ValueString() != "" {
				base["azure_resource_group"] = cm.model.Azure.ResourceGroup.ValueString()
			}
			if !cm.model.Azure.TenantID.IsNull() && !cm.model.Azure.TenantID.IsUnknown() && cm.model.Azure.TenantID.ValueString() != "" {
				base["azure_tenant_id"] = cm.model.Azure.TenantID.ValueString()
			}
			if !cm.model.Azure.ClientID.IsNull() && !cm.model.Azure.ClientID.IsUnknown() && cm.model.Azure.ClientID.ValueString() != "" {
				base["azure_client_id"] = cm.model.Azure.ClientID.ValueString()
			}
			if !cm.model.Azure.ClientSecret.IsNull() && !cm.model.Azure.ClientSecret.IsUnknown() && cm.model.Azure.ClientSecret.ValueString() != "" {
				base["azure_client_secret"] = cm.model.Azure.ClientSecret.ValueString()
			}
		}
	}

	cm.defaultAttrsCache.Store(pType, base)
	return base
}

// DefaultAttrs returns provider-scoped credentials, default regions, and mock mode settings for adapter requests
func (cm *ClientManager) DefaultAttrs(providerType string) map[string]interface{} {
	pType := strings.ToLower(strings.TrimSpace(providerType))
	base := cm.computeStaticDefaultAttrs(pType)

	attrs := make(map[string]interface{}, len(base)+1)
	for k, v := range base {
		attrs[k] = v
	}
	if cm.IsMockMode() {
		attrs["mock_mode"] = true
	}
	return attrs
}

// GetAWSConfig lazily initializes and returns the AWS SDK v2 Config
func (cm *ClientManager) GetAWSConfig(ctx context.Context) (*aws.Config, error) {
	var err error
	cm.awsInitOnce.Do(func() {
		awsRegion := "us-east-1"
		if cm.model.AWS != nil && !cm.model.AWS.Region.IsNull() && !cm.model.AWS.Region.IsUnknown() && cm.model.AWS.Region.ValueString() != "" {
			awsRegion = cm.model.AWS.Region.ValueString()
		}

		opts := []func(*config.LoadOptions) error{
			config.WithRegion(awsRegion),
		}
		if cm.model.AWS != nil {
			if !cm.model.AWS.Profile.IsNull() && !cm.model.AWS.Profile.IsUnknown() && cm.model.AWS.Profile.ValueString() != "" {
				opts = append(opts, config.WithSharedConfigProfile(cm.model.AWS.Profile.ValueString()))
			}
			if !cm.model.AWS.AccessKey.IsNull() && !cm.model.AWS.AccessKey.IsUnknown() && cm.model.AWS.AccessKey.ValueString() != "" &&
				!cm.model.AWS.SecretKey.IsNull() && !cm.model.AWS.SecretKey.IsUnknown() && cm.model.AWS.SecretKey.ValueString() != "" {
				opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					cm.model.AWS.AccessKey.ValueString(),
					cm.model.AWS.SecretKey.ValueString(),
					"",
				)))
			}
		}

		cfg, loadErr := config.LoadDefaultConfig(ctx, opts...)
		err = loadErr
		if err == nil {
			cm.awsConfig = &cfg
		}
	})

	return cm.awsConfig, err
}

// GetGCPConfig lazily initializes and returns the GCP client settings
func (cm *ClientManager) GetGCPConfig(ctx context.Context) (*GCPClientConfig, error) {
	cm.gcpInitOnce.Do(func() {
		gcpProject := "default-project"
		gcpRegion := "us-central1"
		gcpCreds := ""

		if cm.model.GCP != nil {
			if !cm.model.GCP.Project.IsNull() && !cm.model.GCP.Project.IsUnknown() && cm.model.GCP.Project.ValueString() != "" {
				gcpProject = cm.model.GCP.Project.ValueString()
			}
			if !cm.model.GCP.Region.IsNull() && !cm.model.GCP.Region.IsUnknown() && cm.model.GCP.Region.ValueString() != "" {
				gcpRegion = cm.model.GCP.Region.ValueString()
			}
			if !cm.model.GCP.Credentials.IsNull() && !cm.model.GCP.Credentials.IsUnknown() && cm.model.GCP.Credentials.ValueString() != "" {
				gcpCreds = cm.model.GCP.Credentials.ValueString()
			}
		}

		cm.gcpConfig = &GCPClientConfig{
			Project:     gcpProject,
			Region:      gcpRegion,
			Credentials: gcpCreds,
		}
	})

	return cm.gcpConfig, nil
}

// GetAzureConfig lazily initializes and returns the Azure ARM client settings
func (cm *ClientManager) GetAzureConfig(ctx context.Context) (*AzureClientConfig, error) {
	cm.azureInitOnce.Do(func() {
		azureSubID := "00000000-0000-0000-0000-000000000000"
		azureRG := "default-rg"
		var tenantID, clientID, clientSecret string

		if cm.model.Azure != nil {
			if !cm.model.Azure.SubscriptionID.IsNull() && !cm.model.Azure.SubscriptionID.IsUnknown() && cm.model.Azure.SubscriptionID.ValueString() != "" {
				azureSubID = cm.model.Azure.SubscriptionID.ValueString()
			}
			if !cm.model.Azure.ResourceGroup.IsNull() && !cm.model.Azure.ResourceGroup.IsUnknown() && cm.model.Azure.ResourceGroup.ValueString() != "" {
				azureRG = cm.model.Azure.ResourceGroup.ValueString()
			}
			if !cm.model.Azure.TenantID.IsNull() && !cm.model.Azure.TenantID.IsUnknown() {
				tenantID = cm.model.Azure.TenantID.ValueString()
			}
			if !cm.model.Azure.ClientID.IsNull() && !cm.model.Azure.ClientID.IsUnknown() {
				clientID = cm.model.Azure.ClientID.ValueString()
			}
			if !cm.model.Azure.ClientSecret.IsNull() && !cm.model.Azure.ClientSecret.IsUnknown() {
				clientSecret = cm.model.Azure.ClientSecret.ValueString()
			}
		}

		cm.azureConfig = &AzureClientConfig{
			SubscriptionID: azureSubID,
			TenantID:       tenantID,
			ClientID:       clientID,
			ClientSecret:   clientSecret,
			ResourceGroup:  azureRG,
		}
	})

	return cm.azureConfig, nil
}
