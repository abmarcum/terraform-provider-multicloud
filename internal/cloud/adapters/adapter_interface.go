package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/aws"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/azure"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/gcp"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/resiliency"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
)

type ResourceRequest = common.ResourceRequest
type ResourceResponse = common.ResourceResponse
type ClientConfigProvider = common.ClientConfigProvider

var ErrNotFound = common.ErrNotFound

type CloudAdapter interface {
	CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	DeleteResource(ctx context.Context, req ResourceRequest) error
}

type AWSAdapter = aws.AWSAdapter
type GCPAdapter = gcp.GCPAdapter
type AzureAdapter = azure.AzureAdapter

var (
	awsAdapterInstance   = &aws.AWSAdapter{}
	gcpAdapterInstance   = &gcp.GCPAdapter{}
	azureAdapterInstance = &azure.AzureAdapter{}
)

func validatePreApplySecurity(providerType string, resourceType string, resourceName string, extraAttrs map[string]interface{}) error {
	if extraAttrs == nil {
		return nil
	}
	if violations := security.ValidatePolicy(resourceType, resourceName, extraAttrs); len(violations) > 0 {
		return fmt.Errorf("security policy violation (%s): %s", violations[0].RuleName, violations[0].Message)
	}
	for k, v := range extraAttrs {
		// Skip intentional provider authentication keys passed by ClientManager
		if k == "aws_access_key" || k == "aws_secret_key" || k == "gcp_credentials" || k == "azure_bearer_token" || k == "azure_client_secret" {
			continue
		}
		if strVal, ok := v.(string); ok && strVal != "" {
			if leaks := security.ScanForSecretLeaks(providerType, resourceName, fmt.Sprintf("%s = %s", k, strVal)); len(leaks) > 0 {
				return fmt.Errorf("pre-apply secret leak blocked (%s): %s", leaks[0].SecretType, leaks[0].Message)
			}
		}
	}
	return nil
}

func CreateCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	cleanName := sanitizer.SanitizeResourceName(resourceName, providerType, resourceType)
	if err := validatePreApplySecurity(providerType, resourceType, cleanName, extraAttrs); err != nil {
		return ResourceResponse{}, err
	}
	adapter := getAdapter(providerType)
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	}
	return resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.CreateResource(ctx, req)
	})
}

func ReadCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string) (ResourceResponse, error) {
	return ReadCloudResourceWithAttrs(ctx, providerType, resourceType, resourceName, region, nil)
}

func ReadCloudResourceWithAttrs(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	cleanName := sanitizer.SanitizeResourceName(resourceName, providerType, resourceType)
	adapter := getAdapter(providerType)
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	}
	return resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.ReadResource(ctx, req)
	})
}

func UpdateCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	cleanName := sanitizer.SanitizeResourceName(resourceName, providerType, resourceType)
	if err := validatePreApplySecurity(providerType, resourceType, cleanName, extraAttrs); err != nil {
		return ResourceResponse{}, err
	}
	adapter := getAdapter(providerType)
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	}
	return resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.UpdateResource(ctx, req)
	})
}

func DeleteCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string) error {
	return DeleteCloudResourceWithAttrs(ctx, providerType, resourceType, resourceName, region, nil)
}

func DeleteCloudResourceWithAttrs(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) error {
	cleanName := sanitizer.SanitizeResourceName(resourceName, providerType, resourceType)
	adapter := getAdapter(providerType)
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	}
	_, err := resiliency.ExecuteWithRetry(ctx, func() (struct{}, error) {
		return struct{}{}, adapter.DeleteResource(ctx, req)
	})
	return err
}

func getAdapter(providerType string) CloudAdapter {
	switch strings.ToLower(providerType) {
	case "aws":
		return awsAdapterInstance
	case "azure":
		return azureAdapterInstance
	default:
		return gcpAdapterInstance
	}
}
