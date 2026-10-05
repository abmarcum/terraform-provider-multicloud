package gcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GCPAdapter struct{}

var (
	defaultTokenSourceOnce sync.Once
	defaultTokenSource     oauth2.TokenSource
	defaultTokenSourceErr  error
)

func getGCPAccessToken(ctx context.Context, req common.ResourceRequest) (string, error) {
	if req.Attributes != nil {
		if credJSON, ok := req.Attributes["gcp_credentials"].(string); ok && strings.TrimSpace(credJSON) != "" {
			raw := []byte(credJSON)
			if !strings.HasPrefix(strings.TrimSpace(credJSON), "{") {
				/* #nosec G304 */
				if fileBytes, err := os.ReadFile(filepath.Clean(credJSON)); err == nil {
					raw = fileBytes
				}
			}
			creds, err := google.CredentialsFromJSON(ctx, raw, "https://www.googleapis.com/auth/cloud-platform")
			if err != nil {
				return "", fmt.Errorf("GCP credentials parse error: %w", err)
			}
			tok, err := creds.TokenSource.Token()
			if err != nil {
				return "", fmt.Errorf("GCP OAuth2 token exchange failed: %w", err)
			}
			return tok.AccessToken, nil
		}
	}

	defaultTokenSourceOnce.Do(func() {
		defaultTokenSource, defaultTokenSourceErr = google.DefaultTokenSource(context.Background(), "https://www.googleapis.com/auth/cloud-platform")
	})
	if defaultTokenSourceErr != nil || defaultTokenSource == nil {
		return "", fmt.Errorf("GCP authentication failed: %w", defaultTokenSourceErr)
	}
	token, err := defaultTokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("GCP OAuth2 token retrieval failed: %w", err)
	}
	return token.AccessToken, nil
}

func getGCPIntelMachineType(sizeTier string, extraAttrs map[string]interface{}) string {
	if extraAttrs != nil {
		if inst, ok := extraAttrs["instance_type"].(string); ok && inst != "" {
			return inst
		}
		if inst, ok := extraAttrs["gcp_machine_type"].(string); ok && inst != "" {
			return inst
		}
	}
	switch strings.ToLower(sizeTier) {
	case "large":
		return "n2-standard-4"
	case "medium":
		return "n2-standard-2"
	default:
		return "n2-standard-2"
	}
}

func getGCPServiceEndpoint(project string, region string, resType string, name string, extraAttrs map[string]interface{}) (string, []byte) {
	var endpoint string
	var payload []byte

	escProject := url.PathEscape(project)
	escQueryProject := url.QueryEscape(project)
	escRegion := url.PathEscape(region)
	escName := url.PathEscape(name)
	escQueryName := url.QueryEscape(name)

	switch resType {
	case "storage_bucket", "storage_inventory_report":
		endpoint = fmt.Sprintf("https://storage.googleapis.com/storage/v1/b?project=%s", escQueryProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "location": region})
	case "virtual_machine", "bastion_host":
		sizeTier, _ := extraAttrs["size_tier"].(string)
		machineType := getGCPIntelMachineType(sizeTier, extraAttrs)
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instances", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{
			"name":        name,
			"machineType": fmt.Sprintf("zones/%s-a/machineTypes/%s", escRegion, url.PathEscape(machineType)),
		})
	case "custom_machine_type":
		vcpus := int64(2)
		mem := int64(4096)
		if v, ok := extraAttrs["vcpus"].(int64); ok && v > 0 {
			vcpus = v
		}
		if m, ok := extraAttrs["memory_mb"].(int64); ok && m > 0 {
			mem = m
		}
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instances", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{
			"name":        name,
			"machineType": fmt.Sprintf("zones/%s-a/machineTypes/custom-%d-%d", escRegion, vcpus, mem),
		})
	case "virtual_network", "vpc_peering":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/networks", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "autoCreateSubnetworks": true})
	case "subnet":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "ipCidrRange": "10.0.1.0/24"})
	case "security_group", "waf_policy":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/firewalls", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "static_ip":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/addresses", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "global_anycast_ip":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/addresses", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "ipVersion": "IPV4"})
	case "nat_gateway", "vpn_gateway", "transit_gateway":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/routers", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "route_table":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/routes", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "destRange": "0.0.0.0/0"})
	case "load_balancer", "private_endpoint":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/forwardingRules", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "cdn_distribution":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/backendServices", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "enableCDN": true})
	case "auto_scaling_group":
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instanceGroupManagers", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "targetSize": 2})
	case "db_instance":
		endpoint = fmt.Sprintf("https://sqladmin.googleapis.com/v1/projects/%s/instances", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "region": region, "databaseVersion": "POSTGRES_15"})
	case "nosql_table":
		endpoint = fmt.Sprintf("https://firestore.googleapis.com/v1/projects/%s/databases?databaseId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"type": "FIRESTORE_NATIVE", "locationId": region})
	case "cache_cluster":
		endpoint = fmt.Sprintf("https://redis.googleapis.com/v1/projects/%s/locations/%s/instances?instanceId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"tier": "STANDARD_HA", "memorySizeGb": 4})
	case "kubernetes_cluster":
		endpoint = fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"cluster": map[string]interface{}{"name": name, "initialNodeCount": 3}})
	case "container_registry":
		endpoint = fmt.Sprintf("https://artifactregistry.googleapis.com/v1/projects/%s/locations/%s/repositories?repositoryId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"format": "DOCKER"})
	case "container_app":
		endpoint = fmt.Sprintf("https://run.googleapis.com/v2/projects/%s/locations/%s/services?serviceId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"ingress": "INGRESS_TRAFFIC_ALL"})
	case "serverless_function", "edge_function":
		endpoint = fmt.Sprintf("https://cloudfunctions.googleapis.com/v2/projects/%s/locations/%s/functions?functionId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "secret", "secret_rotator":
		endpoint = fmt.Sprintf("https://secretmanager.googleapis.com/v1/projects/%s/secrets?secretId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"replication": map[string]interface{}{"automatic": map[string]interface{}{}}})
	case "kms_key", "kms_policy":
		endpoint = fmt.Sprintf("https://cloudkms.googleapis.com/v1/projects/%s/locations/%s/keyRings/default/cryptoKeys?cryptoKeyId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"purpose": "ENCRYPT_DECRYPT"})
	case "iam_role":
		endpoint = fmt.Sprintf("https://iam.googleapis.com/v1/projects/%s/serviceAccounts", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"accountId": name})
	case "identity_federation", "workload_identity_pool":
		endpoint = fmt.Sprintf("https://iam.googleapis.com/v1/projects/%s/locations/global/workloadIdentityPools?workloadIdentityPoolId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"displayName": name})
	case "dns_zone", "dns_record", "dns_health_check", "dns_zone_link", "dns_resolver", "dnssec", "failover_policy":
		endpoint = fmt.Sprintf("https://dns.googleapis.com/dns/v1/projects/%s/managedZones", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": name, "dnsName": name + ".example.com."})
	case "pubsub_topic", "message_queue":
		endpoint = fmt.Sprintf("https://pubsub.googleapis.com/v1/projects/%s/topics/%s", escProject, escName)
		payload = []byte(`{}`)
	case "event_bridge":
		endpoint = fmt.Sprintf("https://eventarc.googleapis.com/v1/projects/%s/locations/%s/triggers?triggerId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "api_gateway", "graphql_api":
		endpoint = fmt.Sprintf("https://apigateway.googleapis.com/v1/projects/%s/locations/global/apis?apiId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"displayName": name})
	case "data_warehouse":
		endpoint = fmt.Sprintf("https://bigquery.googleapis.com/bigquery/v2/projects/%s/datasets", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"datasetReference": map[string]string{"datasetId": name, "projectId": project}})
	case "search_index":
		endpoint = fmt.Sprintf("https://discoveryengine.googleapis.com/v1/projects/%s/locations/global/collections/default_collection/engines?engineId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"displayName": name})
	case "monitoring_dashboard", "metric_alert":
		endpoint = fmt.Sprintf("https://monitoring.googleapis.com/v1/projects/%s/dashboards", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"displayName": name})
	case "log_workspace":
		endpoint = fmt.Sprintf("https://logging.googleapis.com/v2/projects/%s/locations/%s/buckets?bucketId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"retentionDays": 30})
	case "data_sync", "storage_transfer_job":
		endpoint = "https://storagetransfer.googleapis.com/v1/transferJobs"
		payload, _ = json.Marshal(map[string]interface{}{"projectId": project, "description": name, "status": "ENABLED"})
	case "app_config":
		endpoint = fmt.Sprintf("https://runtimeconfig.googleapis.com/v1beta1/projects/%s/configs", escProject)
		payload, _ = json.Marshal(map[string]interface{}{"name": fmt.Sprintf("projects/%s/configs/%s", project, name)})
	case "ai_endpoint", "feature_store":
		endpoint = fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/endpoints", escRegion, escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"displayName": name})
	case "streaming_cluster":
		endpoint = fmt.Sprintf("https://managedkafka.googleapis.com/v1/projects/%s/locations/%s/clusters?clusterId=%s", escProject, escRegion, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	case "security_center":
		endpoint = fmt.Sprintf("https://securitycenter.googleapis.com/v1/projects/%s/MuteConfigs?muteConfigId=%s", escProject, escQueryName)
		payload, _ = json.Marshal(map[string]interface{}{"filter": "severity=\"HIGH\""})
	case "data_pipeline":
		endpoint = fmt.Sprintf("https://dataflow.googleapis.com/v1b3/projects/%s/locations/%s/jobs", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	default:
		endpoint = fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instances", escProject, escRegion)
		payload, _ = json.Marshal(map[string]interface{}{"name": name})
	}

	return endpoint, payload
}

func getGCPDeleteEndpoint(project string, region string, resType string, name string) string {
	escProject := url.PathEscape(project)
	escRegion := url.PathEscape(region)
	escName := url.PathEscape(name)

	switch resType {
	case "storage_bucket", "storage_inventory_report":
		return fmt.Sprintf("https://storage.googleapis.com/storage/v1/b/%s", escName)
	case "virtual_machine", "custom_machine_type", "bastion_host":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instances/%s", escProject, escRegion, escName)
	case "virtual_network", "vpc_peering":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/networks/%s", escProject, escName)
	case "subnet":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s", escProject, escRegion, escName)
	case "security_group", "waf_policy":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/firewalls/%s", escProject, escName)
	case "static_ip":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/regions/%s/addresses/%s", escProject, escRegion, escName)
	case "global_anycast_ip":
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/addresses/%s", escProject, escName)
	case "db_instance":
		return fmt.Sprintf("https://sqladmin.googleapis.com/v1/projects/%s/instances/%s", escProject, escName)
	case "kubernetes_cluster":
		return fmt.Sprintf("https://container.googleapis.com/v1/projects/%s/locations/%s/clusters/%s", escProject, escRegion, escName)
	case "serverless_function", "edge_function":
		return fmt.Sprintf("https://cloudfunctions.googleapis.com/v2/projects/%s/locations/%s/functions/%s", escProject, escRegion, escName)
	case "secret", "secret_rotator":
		return fmt.Sprintf("https://secretmanager.googleapis.com/v1/projects/%s/secrets/%s", escProject, escName)
	case "pubsub_topic", "message_queue":
		return fmt.Sprintf("https://pubsub.googleapis.com/v1/projects/%s/topics/%s", escProject, escName)
	default:
		return fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s-a/instances/%s", escProject, escRegion, escName)
	}
}

func getGCPReadEndpoint(project string, region string, resType string, name string) string {
	return getGCPDeleteEndpoint(project, region, resType, name)
}

func buildGCPMockAttributes(region string, req common.ResourceRequest) map[string]interface{} {
	attrs := map[string]interface{}{
		"region":         region,
		"provider_type":  "gcp",
		"failoverstatus": "PRIMARY_HEALTHY",
		"syncstatus":     "REPLICATION_ACTIVE",
		"apiendpoint":    fmt.Sprintf("https://%s.apigateway.gcp.cloud.goog", req.ResourceName),
		"ip_address":     fmt.Sprintf("198.51.100.%d", len(req.ResourceName)*7%250+1),
		"dns_name":       fmt.Sprintf("%s.anycast.gcp.net", req.ResourceName),
	}
	for k, v := range req.Attributes {
		attrs[k] = v
	}
	return attrs
}

func (g *GCPAdapter) CreateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	project, _ := common.GetGCPProject(req)
	region := common.GetRegion(req.Region, "us-central1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         fmt.Sprintf("projects/%s/locations/%s/%s/%s", project, region, req.ResourceType, req.ResourceName),
			Status:     "RUNNING",
			Attributes: buildGCPMockAttributes(region, req),
		}, nil
	}

	accessToken, err := getGCPAccessToken(ctx, req)
	if err != nil {
		return common.ResourceResponse{}, err
	}

	apiEndpoint, payload := getGCPServiceEndpoint(project, region, req.ResourceType, req.ResourceName, req.Attributes)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create GCP HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("GCP API call failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 && resp.StatusCode != 409 {
		return common.ResourceResponse{}, fmt.Errorf("GCP API error (status %d): %s", resp.StatusCode, common.SanitizeErrorBody(bodyBytes))
	}

	respAttrs := buildGCPMockAttributes(region, req)
	_ = json.Unmarshal(bodyBytes, &respAttrs)

	return common.ResourceResponse{
		ID:         fmt.Sprintf("projects/%s/locations/%s/%s/%s", project, region, req.ResourceType, req.ResourceName),
		Status:     "RUNNING",
		Attributes: respAttrs,
	}, nil
}

func (g *GCPAdapter) ReadResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	project, _ := common.GetGCPProject(req)
	region := common.GetRegion(req.Region, "us-central1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "RUNNING",
			Attributes: buildGCPMockAttributes(region, req),
		}, nil
	}

	accessToken, err := getGCPAccessToken(ctx, req)
	if err != nil {
		return common.ResourceResponse{}, err
	}

	apiEndpoint := getGCPReadEndpoint(project, region, req.ResourceType, req.ResourceName)
	httpReq, err := http.NewRequestWithContext(ctx, "GET", apiEndpoint, nil)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create GCP Read HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("GCP Read API call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return common.ResourceResponse{}, fmt.Errorf("%w: GCP resource %s not found", common.ErrNotFound, req.ResourceName)
	}
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 400 {
		return common.ResourceResponse{}, fmt.Errorf("GCP API error (status %d) reading %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}

	respAttrs := buildGCPMockAttributes(region, req)
	_ = json.Unmarshal(bodyBytes, &respAttrs)

	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "RUNNING",
		Attributes: respAttrs,
	}, nil
}

func (g *GCPAdapter) UpdateResource(ctx context.Context, req common.ResourceRequest) (common.ResourceResponse, error) {
	project, _ := common.GetGCPProject(req)
	region := common.GetRegion(req.Region, "us-central1")
	if common.IsRequestMockMode(req) {
		return common.ResourceResponse{
			ID:         req.ResourceName,
			Status:     "RUNNING",
			Attributes: buildGCPMockAttributes(region, req),
		}, nil
	}

	accessToken, err := getGCPAccessToken(ctx, req)
	if err != nil {
		return common.ResourceResponse{}, err
	}

	apiEndpoint, payload := getGCPServiceEndpoint(project, region, req.ResourceType, req.ResourceName, req.Attributes)
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", apiEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("failed to create GCP Update HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)
	httpReq.Header.Set("Content-Type", "application/json")

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return common.ResourceResponse{}, fmt.Errorf("GCP Update API call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return common.ResourceResponse{}, fmt.Errorf("GCP API error (status %d) updating %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}

	return common.ResourceResponse{
		ID:         req.ResourceName,
		Status:     "RUNNING",
		Attributes: buildGCPMockAttributes(region, req),
	}, nil
}

func (g *GCPAdapter) DeleteResource(ctx context.Context, req common.ResourceRequest) error {
	project, _ := common.GetGCPProject(req)
	region := common.GetRegion(req.Region, "us-central1")
	if common.IsRequestMockMode(req) {
		return nil
	}

	accessToken, err := getGCPAccessToken(ctx, req)
	if err != nil {
		return err
	}

	apiEndpoint := getGCPDeleteEndpoint(project, region, req.ResourceType, req.ResourceName)
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", apiEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create GCP Delete HTTP request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+accessToken)

	/* #nosec G107 G704 */
	resp, err := common.HTTPClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("GCP Delete API call failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("GCP API error (status %d) deleting %s: %s", resp.StatusCode, req.ResourceName, common.SanitizeErrorBody(bodyBytes))
	}
	return nil
}
