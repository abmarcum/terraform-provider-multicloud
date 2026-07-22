package resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAllUnifiedResourcesMetadataAndSchema(t *testing.T) {
	resourceFactories := map[string]func() resource.Resource{
		"multicloud_storage_bucket":      NewStorageBucketResource,
		"multicloud_virtual_machine":     NewVirtualMachineResource,
		"multicloud_virtual_network":     NewVirtualNetworkResource,
		"multicloud_subnet":              NewSubnetResource,
		"multicloud_security_group":      NewSecurityGroupResource,
		"multicloud_static_ip":           NewStaticIPResource,
		"multicloud_nat_gateway":         NewNATGatewayResource,
		"multicloud_route_table":         NewRouteTableResource,
		"multicloud_load_balancer":       NewLoadBalancerResource,
		"multicloud_db_instance":         NewDBInstanceResource,
		"multicloud_nosql_table":         NewNoSQLTableResource,
		"multicloud_kubernetes_cluster":  NewKubernetesClusterResource,
		"multicloud_container_registry":  NewContainerRegistryResource,
		"multicloud_serverless_function": NewServerlessFunctionResource,
		"multicloud_secret":              NewSecretResource,
		"multicloud_kms_key":             NewKMSKeyResource,
		"multicloud_iam_role":            NewIAMRoleResource,
		"multicloud_dns_zone":            NewDNSZoneResource,
		"multicloud_pubsub_topic":        NewPubSubTopicResource,
		"multicloud_message_queue":       NewMessageQueueResource,
		"multicloud_failover_policy":     NewFailoverPolicyResource,
		"multicloud_event_bridge":        NewEventBridgeResource,
		"multicloud_cdn_distribution":    NewCDNDistributionResource,
		"multicloud_cache_cluster":       NewCacheClusterResource,
		"multicloud_api_gateway":         NewAPIGatewayResource,
		"multicloud_data_warehouse":      NewDataWarehouseResource,
		"multicloud_vpn_gateway":         NewVPNGatewayResource,
		"multicloud_search_index":        NewSearchIndexResource,
		"multicloud_auto_scaling_group":  NewAutoScalingGroupResource,
		"multicloud_monitoring_dashboard": NewMonitoringDashboardResource,
		"multicloud_data_sync":           NewDataSyncResource,
		"multicloud_identity_federation": NewIdentityFederationResource,
		"multicloud_secret_rotator":      NewSecretRotatorResource,
	}

	ctx := context.Background()

	for expectedTypeName, factory := range resourceFactories {
		t.Run(expectedTypeName, func(t *testing.T) {
			r := factory()

			// 1. Validate Metadata
			metaReq := resource.MetadataRequest{ProviderTypeName: "multicloud"}
			metaResp := &resource.MetadataResponse{}
			r.Metadata(ctx, metaReq, metaResp)

			if metaResp.TypeName != expectedTypeName {
				t.Errorf("Metadata TypeName = %s; want %s", metaResp.TypeName, expectedTypeName)
			}

			// 2. Validate Schema
			schemaReq := resource.SchemaRequest{}
			schemaResp := &resource.SchemaResponse{}
			r.Schema(ctx, schemaReq, schemaResp)

			if schemaResp.Schema.Attributes == nil {
				t.Fatalf("Schema attributes for %s is nil", expectedTypeName)
			}

			// Ensure mandatory 'id' attribute exists
			if _, ok := schemaResp.Schema.Attributes["id"]; !ok {
				t.Errorf("Resource %s missing required 'id' attribute in schema", expectedTypeName)
			}
		})
	}
}
