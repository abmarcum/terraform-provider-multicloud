package adapters

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/aws"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/azure"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/gcp"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/resiliency"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/telemetry"
)

type ResourceRequest = common.ResourceRequest
type ResourceResponse = common.ResourceResponse
type AWSAdapter = aws.AWSAdapter
type GCPAdapter = gcp.GCPAdapter
type AzureAdapter = azure.AzureAdapter
type ClientConfigProvider = common.ClientConfigProvider

var (
	ErrNotFound              = common.ErrNotFound
	IsSensitiveOrInternalKey = common.IsSensitiveOrInternalKey
)

type CloudAdapter interface {
	CreateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	ReadResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	UpdateResource(ctx context.Context, req ResourceRequest) (ResourceResponse, error)
	DeleteResource(ctx context.Context, req ResourceRequest) error
}

var (
	awsAdapterInstance   = &aws.AWSAdapter{}
	gcpAdapterInstance   = &gcp.GCPAdapter{}
	azureAdapterInstance = &azure.AzureAdapter{}
)

func GetAdapter(provider string) (CloudAdapter, error) {
	switch {
	case strings.EqualFold(provider, "aws"):
		return awsAdapterInstance, nil
	case strings.EqualFold(provider, "gcp"):
		return gcpAdapterInstance, nil
	case strings.EqualFold(provider, "azure"):
		return azureAdapterInstance, nil
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", provider)
	}
}

func validatePreApplySecurity(provider, resType, cleanName string, extraAttrs map[string]interface{}) error {
	violations := security.ValidatePolicy(resType, cleanName, extraAttrs)
	if len(violations) > 0 {
		return fmt.Errorf("security policy violation [%s]: %s", violations[0].RuleName, violations[0].Message)
	}

	if os.Getenv("OPA_POLICY_PATH") != "" {
		opaRes := security.EvaluateOPARegoPolicyWithAttrs(provider, resType, cleanName, "opa_policy_path", extraAttrs)
		if !opaRes.Passed {
			return fmt.Errorf("OPA policy violation: %s", opaRes.Violation)
		}
	}

	for k, v := range extraAttrs {
		if common.IsSensitiveOrInternalKey(k) {
			continue
		}
		if strVal, ok := v.(string); ok && strVal != "" {
			findings := security.ScanForSecretLeaks(provider, cleanName, strVal)
			if len(findings) > 0 {
				return fmt.Errorf("pre-apply secret scanner blocked provisioning: %s", findings[0].Message)
			}
		}
	}
	return nil
}

func CreateCloudResource(ctx context.Context, provider, resType, name, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	start := time.Now()
	adapter, err := GetAdapter(provider)
	if err != nil {
		return ResourceResponse{}, err
	}

	cleanName := sanitizer.SanitizeResourceName(name, provider, resType)
	cleanRegion := sanitizer.SanitizeCloudIdentifier(region, "")
	if err := validatePreApplySecurity(provider, resType, cleanName, extraAttrs); err != nil {
		return ResourceResponse{}, err
	}

	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resType,
		ProviderType: provider,
		Region:       cleanRegion,
		Attributes:   extraAttrs,
	}
	resp, err := resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.CreateResource(ctx, req)
	})
	_, _ = telemetry.DefaultExporter.RecordEvent("PROVISION", provider, cleanName, time.Since(start), map[string]interface{}{
		"resource_type": resType,
		"success":       err == nil,
	})
	return resp, err
}

func ReadCloudResource(ctx context.Context, provider, resType, name, region string) (ResourceResponse, error) {
	return ReadCloudResourceWithAttrs(ctx, provider, resType, name, region, nil)
}

func ReadCloudResourceWithAttrs(ctx context.Context, provider, resType, name, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	start := time.Now()
	adapter, err := GetAdapter(provider)
	if err != nil {
		return ResourceResponse{}, err
	}

	cleanName := sanitizer.SanitizeResourceName(name, provider, resType)
	cleanRegion := sanitizer.SanitizeCloudIdentifier(region, "")
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resType,
		ProviderType: provider,
		Region:       cleanRegion,
		Attributes:   extraAttrs,
	}
	resp, err := resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.ReadResource(ctx, req)
	})
	_, _ = telemetry.DefaultExporter.RecordEvent("READ", provider, cleanName, time.Since(start), map[string]interface{}{
		"resource_type": resType,
		"success":       err == nil,
	})
	return resp, err
}

func UpdateCloudResource(ctx context.Context, provider, resType, name, region string, extraAttrs map[string]interface{}) (ResourceResponse, error) {
	start := time.Now()
	adapter, err := GetAdapter(provider)
	if err != nil {
		return ResourceResponse{}, err
	}

	cleanName := sanitizer.SanitizeResourceName(name, provider, resType)
	cleanRegion := sanitizer.SanitizeCloudIdentifier(region, "")
	if err := validatePreApplySecurity(provider, resType, cleanName, extraAttrs); err != nil {
		return ResourceResponse{}, err
	}

	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resType,
		ProviderType: provider,
		Region:       cleanRegion,
		Attributes:   extraAttrs,
	}
	resp, err := resiliency.ExecuteWithRetry(ctx, func() (ResourceResponse, error) {
		return adapter.UpdateResource(ctx, req)
	})
	_, _ = telemetry.DefaultExporter.RecordEvent("UPDATE", provider, cleanName, time.Since(start), map[string]interface{}{
		"resource_type": resType,
		"success":       err == nil,
	})
	return resp, err
}

func DeleteCloudResource(ctx context.Context, provider, resType, name, region string) error {
	return DeleteCloudResourceWithAttrs(ctx, provider, resType, name, region, nil)
}

func DeleteCloudResourceWithAttrs(ctx context.Context, provider, resType, name, region string, extraAttrs map[string]interface{}) error {
	start := time.Now()
	adapter, err := GetAdapter(provider)
	if err != nil {
		return err
	}

	cleanName := sanitizer.SanitizeResourceName(name, provider, resType)
	cleanRegion := sanitizer.SanitizeCloudIdentifier(region, "")
	req := ResourceRequest{
		ResourceName: cleanName,
		ResourceType: resType,
		ProviderType: provider,
		Region:       cleanRegion,
		Attributes:   extraAttrs,
	}
	_, err = resiliency.ExecuteWithRetry(ctx, func() (bool, error) {
		if delErr := adapter.DeleteResource(ctx, req); delErr != nil {
			return false, delErr
		}
		return true, nil
	})
	_, _ = telemetry.DefaultExporter.RecordEvent("DELETE", provider, cleanName, time.Since(start), map[string]interface{}{
		"resource_type": resType,
		"success":       err == nil,
	})
	return err
}
