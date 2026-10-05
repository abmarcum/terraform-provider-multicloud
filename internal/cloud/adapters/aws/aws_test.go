package aws

import (
	"context"
	"strings"
	"testing"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters/common"
)

func TestAWSServiceEndpointMappingsAllResources(t *testing.T) {
	resourceTypes := []string{
		"storage_bucket", "virtual_machine", "virtual_network", "subnet", "security_group",
		"static_ip", "nat_gateway", "route_table", "load_balancer", "db_instance",
		"nosql_table", "kubernetes_cluster", "container_registry", "serverless_function",
		"secret", "kms_key", "iam_role", "dns_zone", "pubsub_topic", "message_queue",
		"failover_policy", "event_bridge", "cdn_distribution", "cache_cluster", "api_gateway",
		"data_warehouse", "vpn_gateway", "search_index", "auto_scaling_group",
		"monitoring_dashboard", "data_sync", "identity_federation", "secret_rotator",
		"container_app", "bastion_host", "waf_policy", "vpc_peering", "app_config",
		"ai_endpoint", "streaming_cluster", "metric_alert", "log_workspace", "graphql_api",
		"dns_record", "dns_health_check", "dns_zone_link", "dns_resolver", "dnssec",
		"transit_gateway", "private_endpoint", "security_center", "kms_policy",
		"data_pipeline", "feature_store", "edge_function", "global_anycast_ip",
		"custom_machine_type", "storage_inventory_report", "workload_identity_pool",
		"storage_transfer_job", "block_volume", "shared_filesystem", "tls_certificate",
		"workflow", "batch_compute", "backup_vault", "distributed_tracing", "budget_alert",
		"vector_index", "service_mesh",
	}

	for _, rt := range resourceTypes {
		svc := awsServiceNameForResource(rt)
		if svc == "" {
			t.Errorf("awsServiceNameForResource(%s) returned empty service name", rt)
		}
		ep, method, _ := getAWSServiceEndpoint("us-west-2", rt, "test-res")
		if !strings.HasPrefix(ep, "https://") {
			t.Errorf("getAWSServiceEndpoint(%s) returned invalid endpoint: %s", rt, ep)
		}
		if method == "" {
			t.Errorf("getAWSServiceEndpoint(%s) returned empty HTTP method", rt)
		}
		readEp, readMethod := getAWSReadEndpoint("us-west-2", rt, "test-res")
		if !strings.HasPrefix(readEp, "https://") || readMethod == "" {
			t.Errorf("getAWSReadEndpoint(%s) returned invalid endpoint/method: %s %s", rt, readMethod, readEp)
		}
		delEp, delMethod, _ := getAWSDeleteEndpoint("us-west-2", rt, "test-res")
		if !strings.HasPrefix(delEp, "https://") || delMethod == "" {
			t.Errorf("getAWSDeleteEndpoint(%s) returned invalid endpoint/method: %s %s", rt, delMethod, delEp)
		}
	}
}

func TestAWSAdapterMockLifecycle(t *testing.T) {
	t.Setenv("MULTICLOUD_MOCK_MODE", "true")
	ctx := context.Background()
	adapter := &AWSAdapter{}
	req := common.ResourceRequest{
		ResourceName: "unit-test-res",
		ResourceType: "global_anycast_ip",
		ProviderType: "aws",
		Region:       "us-east-1",
	}

	created, err := adapter.CreateResource(ctx, req)
	if err != nil || created.Status != "ACTIVE" {
		t.Fatalf("CreateResource failed: %v", err)
	}
	if _, ok := created.Attributes["ip_address"].(string); !ok {
		t.Errorf("expected ip_address in AWS mock attributes")
	}

	readResp, err := adapter.ReadResource(ctx, req)
	if err != nil || readResp.Status != "ACTIVE" {
		t.Fatalf("ReadResource failed: %v", err)
	}

	updResp, err := adapter.UpdateResource(ctx, req)
	if err != nil || updResp.Status != "ACTIVE" {
		t.Fatalf("UpdateResource failed: %v", err)
	}

	if err := adapter.DeleteResource(ctx, req); err != nil {
		t.Fatalf("DeleteResource failed: %v", err)
	}
}
