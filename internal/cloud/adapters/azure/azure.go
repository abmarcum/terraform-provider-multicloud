package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
)

type AzureAdapter struct{}

type cachedAzureToken struct {
	token     string
	expiresAt time.Time
}

var (
	azureTokenMu    sync.RWMutex
	azureTokenCache = make(map[string]cachedAzureToken)
)

func getAzureSubscriptionID(req common.ResourceRequest) string {
	sub := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if sub == "" && req.Attributes != nil {
		if s, ok := req.Attributes["azure_subscription_id"].(string); ok && s != "" {
			sub = s
		}
	}
	if sub == "" {
		sub = "00000000-0000-0000-0000-000000000000"
	}
	return url.PathEscape(sub)
}

func getAzureResourceGroup(req common.ResourceRequest) string {
	rg := os.Getenv("AZURE_RESOURCE_GROUP")
	if rg == "" && req.Attributes != nil {
		if r, ok := req.Attributes["azure_resource_group"].(string); ok && r != "" {
			rg = r
		}
	}
	if rg == "" {
		rg = "multicloud-rg"
	}
	return url.PathEscape(rg)
}

func getAzureBearerToken(ctx context.Context, req common.ResourceRequest) string {
	token := os.Getenv("AZURE_BEARER_TOKEN")
	if token != "" {
		return token
	}
	if req.Attributes != nil {
		if t, ok := req.Attributes["azure_bearer_token"].(string); ok && t != "" {
			return t
		}
	}

	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if req.Attributes != nil {
		if tenantID == "" {
			tenantID, _ = req.Attributes["azure_tenant_id"].(string)
		}
		if clientID == "" {
			clientID, _ = req.Attributes["azure_client_id"].(string)
		}
		if clientSecret == "" {
			clientSecret, _ = req.Attributes["azure_client_secret"].(string)
		}
	}

	if tenantID != "" && clientID != "" && clientSecret != "" {
		cacheKey := tenantID + ":" + clientID
		azureTokenMu.RLock()
		if entry, ok := azureTokenCache[cacheKey]; ok && time.Now().Before(entry.expiresAt) {
			azureTokenMu.RUnlock()
			return entry.token
		}
		azureTokenMu.RUnlock()

		tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", url.PathEscape(tenantID))
		form := url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {clientID},
			"client_secret": {clientSecret},
			"scope":         {"https://management.azure.com/.default"},
		}
		/* #nosec G107 G704 */
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
		if err == nil {
			httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			/* #nosec G107 G704 */
			resp, err := common.HTTPClient.Do(httpReq)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					var tokResp struct {
						AccessToken string `json:"access_token"`
						ExpiresIn   int64  `json:"expires_in"`
					}
					if err := json.NewDecoder(io.LimitReader(resp.Body, 65536)).Decode(&tokResp); err == nil && tokResp.AccessToken != "" {
						ttl := 50 * time.Minute
						if tokResp.ExpiresIn > 120 {
							ttl = time.Duration(tokResp.ExpiresIn-60) * time.Second
						}
						azureTokenMu.Lock()
						azureTokenCache[cacheKey] = cachedAzureToken{
							token:     tokResp.AccessToken,
							expiresAt: time.Now().Add(ttl),
						}
						azureTokenMu.Unlock()
						return tokResp.AccessToken
					}
				}
			}
		}
	}

	return ""
}

func getAzureIntelVMSKU(sizeTier string, extraAttrs map[string]interface{}) string {
	if extraAttrs != nil {
		if inst, ok := extraAttrs["instance_type"].(string); ok && inst != "" {
			return inst
		}
		if inst, ok := extraAttrs["azure_vm_sku"].(string); ok && inst != "" {
			return inst
		}
	}
	switch strings.ToLower(sizeTier) {
	case "large":
		return "Standard_D4s_v5"
	case "medium":
		return "Standard_D2s_v5"
	default:
		return "Standard_D2s_v5"
	}
}

func getAzureARMResourcePath(resType string, escName string) (string, string) {
	switch resType {
	case "storage_bucket", "storage_inventory_report":
		return fmt.Sprintf("Microsoft.Storage/storageAccounts/%s", escName), "2023-01-01"
	case "virtual_machine", "custom_machine_type":
		return fmt.Sprintf("Microsoft.Compute/virtualMachines/%s", escName), "2023-09-01"
	case "auto_scaling_group":
		return fmt.Sprintf("Microsoft.Compute/virtualMachineScaleSets/%s", escName), "2023-09-01"
	case "virtual_network", "vpc_peering":
		return fmt.Sprintf("Microsoft.Network/virtualNetworks/%s", escName), "2023-05-01"
	case "subnet":
		return fmt.Sprintf("Microsoft.Network/virtualNetworks/default-vnet/subnets/%s", escName), "2023-05-01"
	case "security_group":
		return fmt.Sprintf("Microsoft.Network/networkSecurityGroups/%s", escName), "2023-05-01"
	case "static_ip":
		return fmt.Sprintf("Microsoft.Network/publicIPAddresses/%s", escName), "2023-05-01"
	case "nat_gateway":
		return fmt.Sprintf("Microsoft.Network/natGateways/%s", escName), "2023-05-01"
	case "route_table":
		return fmt.Sprintf("Microsoft.Network/routeTables/%s", escName), "2023-05-01"
	case "load_balancer":
		return fmt.Sprintf("Microsoft.Network/loadBalancers/%s", escName), "2023-05-01"
	case "vpn_gateway", "transit_gateway":
		return fmt.Sprintf("Microsoft.Network/virtualNetworkGateways/%s", escName), "2023-05-01"
	case "private_endpoint":
		return fmt.Sprintf("Microsoft.Network/privateEndpoints/%s", escName), "2023-05-01"
	case "bastion_host":
		return fmt.Sprintf("Microsoft.Network/bastionHosts/%s", escName), "2023-05-01"
	case "waf_policy":
		return fmt.Sprintf("Microsoft.Network/ApplicationGatewayWebApplicationFirewallPolicies/%s", escName), "2023-05-01"
	case "dns_zone", "dns_record", "dnssec":
		return fmt.Sprintf("Microsoft.Network/dnsZones/%s", escName), "2018-05-01"
	case "dns_zone_link":
		return fmt.Sprintf("Microsoft.Network/privateDnsZones/%s", escName), "2020-06-01"
	case "dns_resolver":
		return fmt.Sprintf("Microsoft.Network/dnsResolvers/%s", escName), "2022-07-01"
	case "dns_health_check", "failover_policy", "global_anycast_ip":
		return fmt.Sprintf("Microsoft.Network/trafficManagerProfiles/%s", escName), "2022-04-01"
	case "db_instance":
		return fmt.Sprintf("Microsoft.DBforPostgreSQL/flexibleServers/%s", escName), "2022-12-01"
	case "nosql_table":
		return fmt.Sprintf("Microsoft.DocumentDB/databaseAccounts/%s", escName), "2023-04-15"
	case "cache_cluster":
		return fmt.Sprintf("Microsoft.Cache/redis/%s", escName), "2023-08-01"
	case "kubernetes_cluster":
		return fmt.Sprintf("Microsoft.ContainerService/managedClusters/%s", escName), "2023-08-01"
	case "container_registry":
		return fmt.Sprintf("Microsoft.ContainerRegistry/registries/%s", escName), "2023-07-01"
	case "container_app":
		return fmt.Sprintf("Microsoft.App/containerApps/%s", escName), "2023-05-01"
	case "serverless_function", "edge_function":
		return fmt.Sprintf("Microsoft.Web/sites/%s", escName), "2022-09-01"
	case "secret", "secret_rotator", "kms_key", "kms_policy":
		return fmt.Sprintf("Microsoft.KeyVault/vaults/%s", escName), "2023-02-01"
	case "iam_role", "identity_federation", "workload_identity_pool":
		return fmt.Sprintf("Microsoft.ManagedIdentity/userAssignedIdentities/%s", escName), "2023-01-31"
	case "pubsub_topic", "event_bridge":
		return fmt.Sprintf("Microsoft.EventGrid/topics/%s", escName), "2022-06-15"
	case "message_queue":
		return fmt.Sprintf("Microsoft.ServiceBus/namespaces/%s", escName), "2022-10-01-preview"
	case "cdn_distribution":
		return fmt.Sprintf("Microsoft.Cdn/profiles/%s", escName), "2023-05-01"
	case "api_gateway", "graphql_api":
		return fmt.Sprintf("Microsoft.ApiManagement/service/%s", escName), "2022-08-01"
	case "data_warehouse":
		return fmt.Sprintf("Microsoft.Synapse/workspaces/%s", escName), "2021-06-01"
	case "search_index":
		return fmt.Sprintf("Microsoft.Search/searchServices/%s", escName), "2022-09-01"
	case "monitoring_dashboard":
		return fmt.Sprintf("Microsoft.Portal/dashboards/%s", escName), "2020-09-01-preview"
	case "metric_alert":
		return fmt.Sprintf("Microsoft.Insights/metricAlerts/%s", escName), "2018-03-01"
	case "log_workspace":
		return fmt.Sprintf("Microsoft.OperationalInsights/workspaces/%s", escName), "2022-10-01"
	case "data_sync", "storage_transfer_job":
		return fmt.Sprintf("Microsoft.StorageSync/storageSyncServices/%s", escName), "2022-06-01"
	case "app_config":
		return fmt.Sprintf("Microsoft.AppConfiguration/configurationStores/%s", escName), "2023-03-01"
	case "ai_endpoint", "feature_store":
		return fmt.Sprintf("Microsoft.MachineLearningServices/workspaces/%s", escName), "2023-04-01"
	case "streaming_cluster":
		return fmt.Sprintf("Microsoft.EventHub/namespaces/%s", escName), "2022-10-01-preview"
	case "security_center":
		return fmt.Sprintf("Microsoft.Security/pricings/%s", escName), "2023-01-01"
	case "data_pipeline":
		return fmt.Sprintf("Microsoft.DataFactory/factories/%s", escName), "2018-06-01"
	default:
		return fmt.Sprintf("Microsoft.Compute/virtualMachines/%s", escName), "2023-09-01"
	}
}

func getAzureServiceEndpoint(subID string, rg string, region string, resType string, name string, extraAttrs map[string]interface{}) (string, []byte) {
	escName := url.PathEscape(name)
	armPath, apiVer := getAzureARMResourcePath(resType, escName)
	endpoint := fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/%s?api-version=%s", subID, rg, armPath, apiVer)

	bodyMap := map[string]interface{}{
		"location": region,
	}
	switch resType {
	case "storage_bucket", "storage_inventory_report":
		bodyMap["sku"] = map[string]string{"name": "Standard_LRS"}
		bodyMap["kind"] = "StorageV2"
	case "virtual_machine", "custom_machine_type":
		sizeTier, _ := extraAttrs["size_tier"].(string)
		vmSize := getAzureIntelVMSKU(sizeTier, extraAttrs)
		bodyMap["properties"] = map[string]interface{}{
			"hardwareProfile": map[string]string{"vmSize": vmSize},
		}
	}

	payload, _ := json.Marshal(bodyMap)
	return endpoint, payload
}

func getAzureDeleteEndpoint(subID string, rg string, resType string, name string) string {
	escName := url.PathEscape(name)
	armPath, apiVer := getAzureARMResourcePath(resType, escName)
	return fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/%s?api-version=%s", subID, rg, armPath, apiVer)
}

func getAzureReadEndpoint(subID string, rg string, resType string, name string) string {
	return getAzureDeleteEndpoint(subID, rg, resType, name)
}

func buildAzureMockAttributes(region string, req common.ResourceRequest) map[string]interface{} {
	attrs := map[string]interface{}{
		"region":         region,
		"provider_type":  "azure",
		"failoverstatus": "PRIMARY_HEALTHY",
		"syncstatus":     "REPLICATION_ACTIVE",
		"apiendpoint":    fmt.Sprintf("https://%s.azure-api.net", req.ResourceName),
		"ip_address":     fmt.Sprintf("198.51.100.%d", len(req.ResourceName)*7%250+1),
		"dns_name":       fmt.Sprintf("%s.anycast.azure.net", req.ResourceName),
	}
	for k, v := range req.Attributes {
		attrs[k] = v
	}
	return attrs
}

func (az *AzureAdapter) CreateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	subID := getAzureSubscriptionID(req)
	rg := getAzureResourceGroup(req)
	region := common.GetRegion(req.Region, "eastus")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MultiCloud/%s/%s", subID, rg, req.ResourceType, req.ResourceName),
			Status:     "SUCCEEDED",
			Attributes: buildAzureMockAttributes(region, req),
		}, nil
	}

	bearerToken := getAzureBearerToken(ctx, req)
	if bearerToken == "" {
		return common.ResourceResponse{}, fmt.Errorf("azure authentication failed: AZURE_BEARER_TOKEN or Azure service principal credentials are required for live ARM provisioning of %s", req.ResourceName)
	}

	apiEndpoint, payload := getAzureServiceEndpoint(subID, rg, region, req.ResourceType, req.ResourceName, req.Attributes)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create Azure HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+bearerToken)
	httpReq.Header.Set("Content-Type", "application/json")

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("azure ARM API call failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return common.ResourceResponse{}, fmt.Errorf("azure ARM API error (status %d): %s", resp.StatusCode, common.SanitizeErrorBody(bodyBytes))
	}

	respAttrs := buildAzureMockAttributes(region, req)
	_ = json.Unmarshal(bodyBytes, &respAttrs)

	return common.ResourceResponse{
		ID:         fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MultiCloud/%s/%s", subID, rg, req.ResourceType, req.ResourceName),
		Status:     "SUCCEEDED",
		Attributes: respAttrs,
	}, nil
}

func (az *AzureAdapter) ReadResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	subID := getAzureSubscriptionID(req)
	rg := getAzureResourceGroup(req)
	region := common.GetRegion(req.Region, "eastus")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "SUCCEEDED",
			Attributes: buildAzureMockAttributes(region, req),
		}, nil
	}

	bearerToken := getAzureBearerToken(ctx, req)
	if bearerToken == "" {
		return common.ResourceResponse{}, fmt.Errorf("azure authentication failed: AZURE_BEARER_TOKEN or Azure service principal credentials are required for live ARM read of %s", req.ResourceName)
	}

	apiEndpoint := getAzureReadEndpoint(subID, rg, req.ResourceType, req.ResourceName)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create Azure Read HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+bearerToken)

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("azure Read ARM API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return common.ResourceResponse{}, fmt.Errorf("%w: Azure resource %s not found", common.ErrNotFound, req.ResourceName)
	}
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return common.ResourceResponse{}, fmt.Errorf("azure ARM API error (status %d) reading %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}

	respAttrs := buildAzureMockAttributes(region, req)
	_ = json.Unmarshal(bodyBytes, &respAttrs)

	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "SUCCEEDED",
		Attributes: respAttrs,
	}, nil
}

func (az *AzureAdapter) UpdateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	subID := getAzureSubscriptionID(req)
	rg := getAzureResourceGroup(req)
	region := common.GetRegion(req.Region, "eastus")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "SUCCEEDED",
			Attributes: buildAzureMockAttributes(region, req),
		}, nil
	}

	bearerToken := getAzureBearerToken(ctx, req)
	if bearerToken == "" {
		return common.ResourceResponse{}, fmt.Errorf("azure authentication failed: AZURE_BEARER_TOKEN or Azure service principal credentials are required for live ARM update of %s", req.ResourceName)
	}

	apiEndpoint, payload := getAzureServiceEndpoint(subID, rg, region, req.ResourceType, req.ResourceName, req.Attributes)
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create Azure Update HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+bearerToken)
	httpReq.Header.Set("Content-Type", "application/json")

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("azure Update ARM API call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return common.ResourceResponse{}, fmt.Errorf("azure ARM API error (status %d) updating %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}

	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "SUCCEEDED",
		Attributes: buildAzureMockAttributes(region, req),
	}, nil
}

func (az *AzureAdapter) DeleteResource(ctx context.Context, req common.ResourceRequest) error {
	subID := getAzureSubscriptionID(req)
	rg := getAzureResourceGroup(req)
	if common.IsRequestMockMode(req) {
		return nil
	}

	bearerToken := getAzureBearerToken(ctx, req)
	if bearerToken == "" {
		return fmt.Errorf("azure authentication failed: AZURE_BEARER_TOKEN or Azure service principal credentials are required for live ARM deletion of %s", req.ResourceName)
	}

	apiEndpoint := getAzureDeleteEndpoint(subID, rg, req.ResourceType, req.ResourceName)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create Azure Delete HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+bearerToken)

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("azure Delete ARM API call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("azure ARM API error (status %d) deleting %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}
	return nil
}
