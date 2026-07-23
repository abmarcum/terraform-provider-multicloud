package adapters

import (
	"context"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/aws"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/azure"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/gcp"
)

type ResourceRequest = common.ResourceRequest
type ResourceResponse = common.ResourceResponse

type CloudAdapter interface {
	CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	DeleteResource(ctx context.Context, req ResourceRequest) error
}

type AWSAdapter = aws.AWSAdapter
type GCPAdapter = gcp.GCPAdapter
type AzureAdapter = azure.AzureAdapter

func CreateCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	adapter := getAdapter(providerType)
	return adapter.CreateResource(ctx, ResourceRequest{
		ResourceName: resourceName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	})
}

func ReadCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string) (ResourceResponse, error) {
	adapter := getAdapter(providerType)
	return adapter.ReadResource(ctx, ResourceRequest{
		ResourceName: resourceName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
	})
}

func UpdateCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	adapter := getAdapter(providerType)
	return adapter.UpdateResource(ctx, ResourceRequest{
		ResourceName: resourceName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
		Attributes:   extraAttrs,
	})
}

func DeleteCloudResource(ctx context.Context, providerType string, resourceType string, resourceName string, region string) error {
	adapter := getAdapter(providerType)
	return adapter.DeleteResource(ctx, ResourceRequest{
		ResourceName: resourceName,
		ResourceType: resourceType,
		ProviderType: providerType,
		Region:       region,
	})
}

func getAdapter(providerType string) CloudAdapter {
	switch strings.ToLower(providerType) {
	case "aws":
		return &aws.AWSAdapter{}
	case "azure":
		return &azure.AzureAdapter{}
	default:
		return &gcp.GCPAdapter{}
	}
}
