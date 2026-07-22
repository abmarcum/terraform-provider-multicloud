package provider

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// ClientManager holds cloud SDK instances with thread-safe lazy initialization
type ClientManager struct {
	mu            sync.Mutex
	model         ProviderModel
	awsConfig     *aws.Config
	gcpConfig     *GCPClientConfig
	azureConfig   *AzureClientConfig
	awsInitOnce   sync.Once
	gcpInitOnce   sync.Once
	azureInitOnce sync.Once
}

type GCPClientConfig struct {
	Project string
	Region  string
}

type AzureClientConfig struct {
	SubscriptionID string
	ResourceGroup  string
}

// NewClientManager returns a ClientManager that defers SDK initialization until clients are actually requested
func NewClientManager(ctx context.Context, model ProviderModel) (*ClientManager, error) {
	cm := &ClientManager{
		model: model,
	}
	return cm, nil
}

// GetAWSConfig lazily initializes and returns the AWS SDK v2 Config
func (cm *ClientManager) GetAWSConfig(ctx context.Context) (*aws.Config, error) {
	var err error
	cm.awsInitOnce.Do(func() {
		awsRegion := "us-east-1"
		if cm.model.AWS != nil && !cm.model.AWS.Region.IsNull() && !cm.model.AWS.Region.IsUnknown() && cm.model.AWS.Region.ValueString() != "" {
			awsRegion = cm.model.AWS.Region.ValueString()
		}

		var cfg aws.Config
		if cm.model.AWS != nil &&
			!cm.model.AWS.AccessKey.IsNull() && !cm.model.AWS.AccessKey.IsUnknown() && cm.model.AWS.AccessKey.ValueString() != "" &&
			!cm.model.AWS.SecretKey.IsNull() && !cm.model.AWS.SecretKey.IsUnknown() && cm.model.AWS.SecretKey.ValueString() != "" {
			cfg, err = config.LoadDefaultConfig(ctx,
				config.WithRegion(awsRegion),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
					cm.model.AWS.AccessKey.ValueString(),
					cm.model.AWS.SecretKey.ValueString(),
					"",
				)),
			)
		} else {
			cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(awsRegion))
		}

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

		if cm.model.GCP != nil {
			if !cm.model.GCP.Project.IsNull() && !cm.model.GCP.Project.IsUnknown() && cm.model.GCP.Project.ValueString() != "" {
				gcpProject = cm.model.GCP.Project.ValueString()
			}
			if !cm.model.GCP.Region.IsNull() && !cm.model.GCP.Region.IsUnknown() && cm.model.GCP.Region.ValueString() != "" {
				gcpRegion = cm.model.GCP.Region.ValueString()
			}
		}

		cm.gcpConfig = &GCPClientConfig{
			Project: gcpProject,
			Region:  gcpRegion,
		}
	})

	return cm.gcpConfig, nil
}

// GetAzureConfig lazily initializes and returns the Azure ARM client settings
func (cm *ClientManager) GetAzureConfig(ctx context.Context) (*AzureClientConfig, error) {
	cm.azureInitOnce.Do(func() {
		azureSubID := "00000000-0000-0000-0000-000000000000"
		azureRG := "default-rg"

		if cm.model.Azure != nil {
			if !cm.model.Azure.SubscriptionID.IsNull() && !cm.model.Azure.SubscriptionID.IsUnknown() && cm.model.Azure.SubscriptionID.ValueString() != "" {
				azureSubID = cm.model.Azure.SubscriptionID.ValueString()
			}
			if !cm.model.Azure.ResourceGroup.IsNull() && !cm.model.Azure.ResourceGroup.IsUnknown() && cm.model.Azure.ResourceGroup.ValueString() != "" {
				azureRG = cm.model.Azure.ResourceGroup.ValueString()
			}
		}

		cm.azureConfig = &AzureClientConfig{
			SubscriptionID: azureSubID,
			ResourceGroup:  azureRG,
		}
	})

	return cm.azureConfig, nil
}
