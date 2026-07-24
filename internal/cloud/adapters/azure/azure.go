package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
)

type AzureAdapter struct{}

func getAzureResourceType(resType string) (string, string) {
	switch resType {
	case "storage_bucket":
		return "Microsoft.Storage/storageAccounts", "2021-04-01"
	case "virtual_network":
		return "Microsoft.Network/virtualNetworks", "2021-02-01"
	case "subnet":
		return "Microsoft.Network/virtualNetworks/subnets", "2021-02-01"
	case "security_group":
		return "Microsoft.Network/networkSecurityGroups", "2021-02-01"
	case "db_instance":
		return "Microsoft.Sql/servers", "2021-02-01-preview"
	case "secret", "secret_rotator":
		return "Microsoft.KeyVault/vaults/secrets", "2021-06-01-preview"
	case "serverless_function":
		return "Microsoft.Web/sites", "2021-02-01"
	case "kubernetes_cluster":
		return "Microsoft.ContainerService/managedClusters", "2021-05-01"
	case "cache_cluster":
		return "Microsoft.Cache/Redis", "2020-06-01"
	case "pubsub_topic":
		return "Microsoft.EventHub/namespaces/eventhubs", "2021-06-01-preview"
	case "bastion_host", "virtual_machine":
		return "Microsoft.Compute/virtualMachines", "2021-03-01"
	default:
		return "Microsoft.Resources/deployments", "2021-04-01"
	}
}

func getAzureSubscriptionID() string {
	sub := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if sub == "" {
		sub = "unconfigured-subscription-id"
	}
	return url.PathEscape(sub)
}

func getAzureResourceURI(subID string, rg string, resType string, name string) string {
	escSubID := url.PathEscape(subID)
	escRG := url.PathEscape(rg)
	escName := url.PathEscape(name)
	azType, apiVer := getAzureResourceType(resType)

	if azType == "Microsoft.Resources/deployments" {
		return fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s?api-version=%s",
			escSubID, escRG, escName, apiVer)
	}
	return fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/%s/%s?api-version=%s",
		escSubID, escRG, azType, escName, apiVer)
}

func getAzureVMHardwareSKU(sizeTier string, extraAttrs map[string]interface{}) string {
	if extraAttrs != nil {
		if sku, ok := extraAttrs["azure_vm_sku"].(string); ok && sku != "" {
			return sku
		}
		if arch, ok := extraAttrs["azure_hardware_architecture"].(string); ok && strings.ToLower(arch) == "intel" {
			switch strings.ToLower(sizeTier) {
			case "large":
				return "Standard_D8s_v5" // Intel Xeon Platinum 8370C (Ice Lake)
			case "medium":
				return "Standard_D4s_v5" // Intel Xeon Platinum 8370C (Ice Lake)
			default:
				return "Standard_D2s_v5" // Intel Xeon Platinum 8370C (Ice Lake)
			}
		}
		if confidential, ok := extraAttrs["intel_sgx_confidential"].(string); ok && confidential == "true" {
			return "Standard_DC2s_v3" // Intel SGX Confidential Compute with Intel Xeon E-2288G
		}
	}

	switch strings.ToLower(sizeTier) {
	case "large":
		return "Standard_D8s_v5" // Intel Xeon Platinum 8370C
	case "medium":
		return "Standard_D4s_v5" // Intel Xeon Platinum 8370C
	default:
		return "Standard_D2s_v5" // Intel Xeon Platinum 8370C
	}
}

func (a *AzureAdapter) CreateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	subscriptionID := getAzureSubscriptionID()
	resourceGroup := "multicloud-rg"
	region := common.GetRegion(req.Region, "eastus")

	if common.IsMockMode() {
		return common.ResourceResponse{
			ID:     fmt.Sprintf("azure/%s/%s/%s", subscriptionID, resourceGroup, req.ResourceName),
			Status: "SUCCEEDED",
		}, nil
	}

	token := os.Getenv("AZURE_BEARER_TOKEN")
	if token == "" {
		return common.ResourceResponse{}, fmt.Errorf("Azure authentication error: AZURE_BEARER_TOKEN is not set")
	}

	azType, apiVer := getAzureResourceType(req.ResourceType)
	armEndpoint := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Resources/deployments/%s?api-version=2021-04-01",
		subscriptionID, url.PathEscape(resourceGroup), url.PathEscape(req.ResourceName))

	sizeTier := "small"
	if s, ok := req.Attributes["size_tier"].(string); ok && s != "" {
		sizeTier = s
	}
	intelVmSku := getAzureVMHardwareSKU(sizeTier, req.Attributes)

	vmProperties := map[string]interface{}{
		"hardwareProfile": map[string]string{
			"vmSize": intelVmSku, // Intel Xeon Platinum / SGX hardware SKU
		},
	}

	armPayload := map[string]interface{}{
		"properties": map[string]interface{}{
			"mode": "Incremental",
			"template": map[string]interface{}{
				"$schema":        "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
				"contentVersion": "1.0.0.0",
				"resources": []map[string]interface{}{
					{
						"type":       azType,
						"apiVersion": apiVer,
						"name":       req.ResourceName,
						"location":   region,
						"properties": vmProperties,
					},
				},
			},
		},
	}
	payload, _ := json.Marshal(armPayload)

	httpReq, err := http.NewRequestWithContext(ctx, "PUT", armEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create Azure HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("Azure ARM API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return common.ResourceResponse{}, fmt.Errorf("Azure ARM API error (status %d) for %s", resp.StatusCode, req.ResourceName)
	}

	return common.ResourceResponse{
		ID:     fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s", subscriptionID, resourceGroup, azType, req.ResourceName),
		Status: "SUCCEEDED",
	}, nil
}

func (a *AzureAdapter) ReadResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	subscriptionID := getAzureSubscriptionID()
	resourceGroup := "multicloud-rg"

	if common.IsMockMode() {
		return common.ResourceResponse{ID: req.ResourceName, Status: "SUCCEEDED"}, nil
	}

	token := os.Getenv("AZURE_BEARER_TOKEN")
	if token == "" {
		return common.ResourceResponse{}, fmt.Errorf("Azure authentication error: AZURE_BEARER_TOKEN is not set")
	}

	armEndpoint := getAzureResourceURI(subscriptionID, resourceGroup, req.ResourceType, req.ResourceName)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", armEndpoint, nil)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create Azure Read HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("Azure Read API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return common.ResourceResponse{}, fmt.Errorf("Azure resource %s read error (status %d)", req.ResourceName, resp.StatusCode)
	}
	return common.ResourceResponse{ID: req.ResourceName, Status: "SUCCEEDED"}, nil
}

func (a *AzureAdapter) UpdateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	return a.CreateResource(ctx, req)
}

func (a *AzureAdapter) DeleteResource(ctx context.Context, req common.ResourceRequest) error {
	subscriptionID := getAzureSubscriptionID()
	resourceGroup := "multicloud-rg"

	armEndpoint := getAzureResourceURI(subscriptionID, resourceGroup, req.ResourceType, req.ResourceName)
	token := os.Getenv("AZURE_BEARER_TOKEN")
	if token != "" {
		httpReq, err := http.NewRequestWithContext(ctx, "DELETE", armEndpoint, nil)
		if err == nil {
			httpReq.Header.Set("Authorization", "Bearer "+token)
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}
	return nil
}
