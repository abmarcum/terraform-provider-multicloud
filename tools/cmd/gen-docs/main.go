package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type ResourceMeta struct {
	Name            string
	Category        string
	Description     string
	AWS             string
	GCP             string
	Azure           string
	Attributes      []string
	BasicExample    string
	AdvancedExample string
	HowItWorks      string
}

func main() {
	fmt.Println("======================================================================")
	fmt.Println("  AUTOMATED TERRAFORM REGISTRY DOCUMENTATION GENERATOR (33 RESOURCES)")
	fmt.Println("======================================================================")

	docsDir := "docs/resources"
	/* #nosec G301 */
	if err := os.MkdirAll(docsDir, 0750); err != nil {
		fmt.Printf("Error creating docs directory: %v\n", err)
		os.Exit(1)
	}

	resources := []ResourceMeta{
		{
			Name:        "storage_bucket",
			Category:    "Storage",
			Description: "Unified Object Storage Bucket supporting AWS S3, GCP Cloud Storage, and Azure Storage Container.",
			AWS:         "aws_s3_bucket",
			GCP:         "google_storage_bucket",
			Azure:       "azurerm_storage_container",
			HowItWorks:  "The `multicloud_storage_bucket` resource provisions object storage containers across AWS, GCP, and Azure. It sanitizes bucket names to meet cloud-specific DNS constraints, enables versioning and server-side encryption pre-apply, and maps extra S3/GCS/ARM settings via `extra_config`.",
			Attributes:  []string{"`bucket_name` (String, Required) Name of the storage bucket.", "`versioning_enabled` (Bool, Optional) Enable object versioning.", "`encryption_enabled` (Bool, Optional) Enable server-side encryption."},
			BasicExample: `resource "multicloud_storage_bucket" "aws_basic" {
  provider_type      = "aws"
  bucket_name        = "my-app-storage-bucket-aws"
  region             = "us-west-2"
  versioning_enabled = true
  encryption_enabled = true
}`,
			AdvancedExample: `resource "multicloud_storage_bucket" "gcp_advanced" {
  provider_type      = "gcp"
  bucket_name        = "my-app-storage-bucket-gcp"
  region             = "us-central1"
  versioning_enabled = true
  encryption_enabled = true

  extra_config = {
    "gcp_storage_class" = "NEARLINE"
    "gcp_location_type" = "dual-region"
  }
}`,
		},
		{
			Name:        "virtual_machine",
			Category:    "Compute",
			Description: "Unified Virtual Machine Compute Instance supporting AWS EC2, GCP Compute Engine, and Azure VMs.",
			AWS:         "aws_instance",
			GCP:         "google_compute_instance",
			Azure:       "azurerm_linux_virtual_machine",
			HowItWorks:  "The `multicloud_virtual_machine` resource maps standardized compute tier profiles ('small', 'medium', 'large') to native instance types (`t3.medium`, `e2-standard-2`, `Standard_B2s`). It manages SSH key deployment, VPC subnet association, and tags across AWS, GCP, and Azure.",
			Attributes:  []string{"`vm_name` (String, Required) Virtual machine instance name.", "`size_tier` (String, Optional) Instance tier ('small', 'medium', 'large').", "`image_id` (String, Optional) Custom image ID.", "`subnet_id` (String, Optional) VPC subnet ID.", "`ssh_public_key` (String, Optional, Sensitive) Public SSH key."},
			BasicExample: `resource "multicloud_virtual_machine" "gcp_basic" {
  provider_type = "gcp"
  vm_name       = "web-server-gcp"
  size_tier     = "medium"
  region        = "us-central1"
}`,
			AdvancedExample: `resource "multicloud_virtual_machine" "aws_advanced" {
  provider_type = "aws"
  vm_name       = "database-host-aws"
  size_tier     = "large"
  region        = "us-west-2"
  subnet_id     = "subnet-0123456789abcdef0"

  extra_config = {
    "aws_ebs_optimized" = "true"
    "aws_monitoring"    = "detailed"
  }
}`,
		},
		{
			Name:        "failover_policy",
			Category:    "Disaster Recovery",
			Description: "Unified Active-Passive Disaster Recovery Failover Policy.",
			AWS:         "Route53 Failover Routing",
			GCP:         "Cloud DNS Failover Policy",
			Azure:       "Azure Traffic Manager",
			HowItWorks:  "The `multicloud_failover_policy` resource configures global active-passive routing between primary and secondary clouds. During runtime, managed cloud DNS edge networks (AWS Route 53, GCP Cloud DNS, Azure Traffic Manager) probe `health_check_url` every 10-30s. If the primary cloud fails, traffic is automatically diverted to the failover cloud without local agent dependencies.",
			Attributes:  []string{"`policy_name` (String, Required) Policy identifier.", "`primary_cloud` (String, Required) Primary cloud provider ('aws', 'gcp', 'azure').", "`failover_cloud` (String, Required) Backup cloud provider.", "`health_check_url` (String, Optional) Endpoint URL.", "`auto_failover` (Bool, Optional) Automated failover status."},
			BasicExample: `resource "multicloud_failover_policy" "basic" {
  policy_name      = "global-app-failover"
  primary_cloud    = "aws"
  failover_cloud   = "gcp"
  health_check_url = "https://app.example.com/healthz"
  auto_failover    = true
}`,
			AdvancedExample: `resource "multicloud_failover_policy" "azure_advanced" {
  policy_name      = "enterprise-multi-cloud-dr"
  primary_cloud    = "azure"
  failover_cloud   = "aws"
  health_check_url = "https://api.example.com/healthz"
  auto_failover    = true

  extra_config = {
    "azure_routing_method" = "Performance"
    "azure_ttl"            = "60"
  }
}`,
		},
		{
			Name:        "data_sync",
			Category:    "Replication",
			Description: "Multi-cloud object storage background data sync.",
			AWS:         "S3 Cross-Region Replication",
			GCP:         "Storage Transfer Service",
			Azure:       "Azure Storage Sync",
			HowItWorks:  "The `multicloud_data_sync` resource orchestrates cloud-native background object replication jobs between source and destination storage buckets across cloud boundaries using S3 Replication, GCP Storage Transfer Service, or Azure Storage Sync.",
			Attributes:  []string{"`sync_name` (String, Required) Task identifier.", "`source_bucket` (String, Optional) Source URI.", "`destination_bucket` (String, Optional) Destination URI.", "`schedule_cron` (String, Optional) Cron schedule."},
			BasicExample: `resource "multicloud_data_sync" "basic" {
  provider_type      = "aws"
  sync_name          = "s3-to-gcs-sync"
  source_bucket      = "s3://my-source-bucket"
  destination_bucket = "gs://my-dest-bucket"
}`,
			AdvancedExample: `resource "multicloud_data_sync" "advanced" {
  provider_type      = "gcp"
  sync_name          = "nightly-analytics-sync"
  source_bucket      = "gs://prod-logs-bucket"
  destination_bucket = "s3://analytics-archive-bucket"
  schedule_cron      = "0 2 * * *"

  extra_config = {
    "gcp_overwrite_objects" = "true"
  }
}`,
		},
		{
			Name:        "identity_federation",
			Category:    "Security",
			Description: "Multi-cloud OIDC workload identity federation.",
			AWS:         "IAM OIDC Provider",
			GCP:         "Workload Identity",
			Azure:       "Entra ID Workload Identity",
			HowItWorks:  "The `multicloud_identity_federation` resource configures passwordless OIDC workload identity trusts across AWS IAM OIDC, GCP Workload Identity, and Azure Entra ID, allowing external workloads (e.g. Kubernetes, GitHub Actions) to authenticate securely without long-lived static secrets.",
			Attributes:  []string{"`federation_name` (String, Required) Federation identifier.", "`issuer_url` (String, Optional) OIDC issuer URL.", "`client_id_list` (List[String], Optional) Target client IDs."},
			BasicExample: `resource "multicloud_identity_federation" "basic" {
  provider_type   = "aws"
  federation_name = "k8s-oidc"
  issuer_url      = "https://container.googleapis.com/v1/projects/my-proj/locations/us-central1/clusters/my-cluster"
}`,
			AdvancedExample: `resource "multicloud_identity_federation" "azure_advanced" {
  provider_type   = "azure"
  federation_name = "github-actions-federation"
  issuer_url      = "https://token.actions.githubusercontent.com"

  extra_config = {
    "azure_audience" = "api://AzureADTokenExchange"
  }
}`,
		},
		{
			Name:        "secret_rotator",
			Category:    "Security",
			Description: "Multi-cloud automated secret key rotator schedule.",
			AWS:         "Secrets Manager Rotation",
			GCP:         "Secret Manager Rotation",
			Azure:       "Key Vault Secret Rotation",
			HowItWorks:  "The `multicloud_secret_rotator` resource configures automated secret rotation schedules in AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault, triggering key generation functions at specified interval days.",
			Attributes:  []string{"`rotator_name` (String, Required) Schedule identifier.", "`secret_id` (String, Optional) Target secret ID.", "`rotation_days` (Int64, Optional) Interval in days."},
			BasicExample: `resource "multicloud_secret_rotator" "basic" {
  provider_type = "aws"
  rotator_name  = "db-key-rotator"
  rotation_days = 90
}`,
			AdvancedExample: `resource "multicloud_secret_rotator" "advanced" {
  provider_type = "azure"
  rotator_name  = "api-key-rotator-azure"
  rotation_days = 30

  extra_config = {
    "azure_notify_before_expiry" = "P7D"
  }
}`,
		},
		{
			Name:        "auto_scaling_group",
			Category:    "Compute",
			Description: "Unified Compute Auto Scaling Group.",
			AWS:         "aws_autoscaling_group",
			GCP:         "google_compute_instance_group_manager",
			Azure:       "azurerm_linux_virtual_machine_scale_set",
			HowItWorks:  "The `multicloud_auto_scaling_group` resource provisions dynamic compute scale groups across AWS ASGs, GCP Instance Group Managers, and Azure VM Scale Sets, managing minimum, maximum, and desired node capacities.",
			Attributes:  []string{"`group_name` (String, Required) Group name.", "`min_size` (Int64, Optional) Minimum capacity.", "`max_size` (Int64, Optional) Maximum capacity.", "`desired_capacity` (Int64, Optional) Desired capacity."},
			BasicExample: `resource "multicloud_auto_scaling_group" "basic" {
  provider_type    = "aws"
  group_name       = "web-asg"
  min_size         = 2
  max_size         = 10
  desired_capacity = 4
}`,
			AdvancedExample: `resource "multicloud_auto_scaling_group" "gcp_advanced" {
  provider_type    = "gcp"
  group_name       = "frontend-mig-gcp"
  min_size         = 3
  max_size         = 15
  desired_capacity = 5

  extra_config = {
    "gcp_target_cpu_utilization" = "0.75"
  }
}`,
		},
		{
			Name:        "serverless_function",
			Category:    "Compute",
			Description: "Unified Serverless Event Function.",
			AWS:         "aws_lambda_function",
			GCP:         "google_cloudfunctions2_function",
			Azure:       "azurerm_linux_function_app",
			HowItWorks:  "The `multicloud_serverless_function` resource provisions event-driven serverless code execution units across AWS Lambda, GCP Cloud Functions v2, and Azure Function Apps, managing runtimes, memory limits, and environment variables.",
			Attributes:  []string{"`function_name` (String, Required) Function name.", "`runtime` (String, Optional) Runtime environment.", "`handler` (String, Optional) Entrypoint handler.", "`memory_size_mb` (Int64, Optional) Memory in MB.", "`timeout_seconds` (Int64, Optional) Timeout in seconds."},
			BasicExample: `resource "multicloud_serverless_function" "basic" {
  provider_type   = "aws"
  function_name   = "image-resizer"
  runtime         = "python3.11"
  handler         = "main.handler"
  memory_size_mb  = 512
}`,
			AdvancedExample: `resource "multicloud_serverless_function" "gcp_advanced" {
  provider_type   = "gcp"
  function_name   = "event-processor-gcp"
  runtime         = "nodejs20"
  handler         = "processEvent"
  memory_size_mb  = 1024
  timeout_seconds = 60

  extra_config = {
    "gcp_max_instance_count" = "50"
  }
}`,
		},
		{
			Name:        "virtual_network",
			Category:    "Network",
			Description: "Unified Virtual Network (VPC / VNet).",
			AWS:         "aws_vpc",
			GCP:         "google_compute_network",
			Azure:       "azurerm_virtual_network",
			HowItWorks:  "The `multicloud_virtual_network` resource provisions isolated virtual private networks across AWS VPCs, GCP VPC Networks, and Azure VNets, configuring CIDR address spaces and network routing boundaries.",
			Attributes:  []string{"`network_name` (String, Required) Network identifier.", "`cidr_block` (String, Optional) IPv4 CIDR range."},
			BasicExample: `resource "multicloud_virtual_network" "basic" {
  provider_type = "aws"
  network_name  = "prod-vpc"
  cidr_block    = "10.0.0.0/16"
}`,
			AdvancedExample: `resource "multicloud_virtual_network" "azure_advanced" {
  provider_type = "azure"
  network_name  = "corp-vnet-azure"
  cidr_block    = "172.16.0.0/12"

  extra_config = {
    "azure_enable_ddos_protection" = "true"
  }
}`,
		},
		{
			Name:        "subnet",
			Category:    "Network",
			Description: "Unified Subnetwork.",
			AWS:         "aws_subnet",
			GCP:         "google_compute_subnetwork",
			Azure:       "azurerm_subnet",
			HowItWorks:  "The `multicloud_subnet` resource provisions subnets within virtual networks across AWS Subnets, GCP Subnetworks, and Azure Subnets, associating IPv4 CIDR blocks and availability zone placement.",
			Attributes:  []string{"`subnet_name` (String, Required) Subnet identifier.", "`vpc_id` (String, Optional) Parent VPC ID.", "`cidr_block` (String, Optional) Subnet CIDR block.", "`availability_zone` (String, Optional) AZ placement."},
			BasicExample: `resource "multicloud_subnet" "basic" {
  provider_type     = "gcp"
  subnet_name       = "public-subnet-1"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-central1-a"
}`,
			AdvancedExample: `resource "multicloud_subnet" "aws_advanced" {
  provider_type     = "aws"
  subnet_name       = "private-subnet-az2"
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-west-2b"

  extra_config = {
    "aws_map_public_ip_on_launch" = "false"
  }
}`,
		},
		{
			Name:        "static_ip",
			Category:    "Network",
			Description: "Unified Static Elastic IP Address.",
			AWS:         "aws_eip",
			GCP:         "google_compute_address",
			Azure:       "azurerm_public_ip",
			HowItWorks:  "The `multicloud_static_ip` resource allocates static public IPv4 addresses across AWS Elastic IPs, GCP Compute Addresses, and Azure Public IPs.",
			Attributes:  []string{"`ip_name` (String, Required) IP name.", "`allocation_type` (String, Optional) 'static' or 'dynamic'."},
			BasicExample: `resource "multicloud_static_ip" "basic" {
  provider_type   = "azure"
  ip_name         = "app-public-ip"
  allocation_type = "static"
}`,
			AdvancedExample: `resource "multicloud_static_ip" "gcp_advanced" {
  provider_type   = "gcp"
  ip_name         = "lb-static-ip-gcp"
  allocation_type = "static"

  extra_config = {
    "gcp_network_tier" = "PREMIUM"
  }
}`,
		},
		{
			Name:        "nat_gateway",
			Category:    "Network",
			Description: "Unified NAT Gateway.",
			AWS:         "aws_nat_gateway",
			GCP:         "google_compute_router_nat",
			Azure:       "azurerm_nat_gateway",
			HowItWorks:  "The `multicloud_nat_gateway` resource provisions outbound Network Address Translation (NAT) gateways across AWS NAT Gateways, GCP Cloud NAT, and Azure NAT Gateways to grant private subnet workloads internet egress.",
			Attributes:  []string{"`gateway_name` (String, Required) Gateway name.", "`subnet_id` (String, Optional) Target Subnet ID.", "`allocation_id` (String, Optional) Static IP ID."},
			BasicExample: `resource "multicloud_nat_gateway" "basic" {
  provider_type = "aws"
  gateway_name  = "main-nat-gw"
}`,
			AdvancedExample: `resource "multicloud_nat_gateway" "advanced" {
  provider_type = "aws"
  gateway_name  = "ha-nat-gw"
  subnet_id     = "subnet-abc12345"

  extra_config = {
    "aws_secondary_allocation_ids" = "eipalloc-98765432"
  }
}`,
		},
		{
			Name:        "route_table",
			Category:    "Network",
			Description: "Unified Routing Table.",
			AWS:         "aws_route_table",
			GCP:         "google_compute_route",
			Azure:       "azurerm_route_table",
			HowItWorks:  "The `multicloud_route_table` resource configures routing rules and gateway associations for virtual networks across AWS Route Tables, GCP Compute Routes, and Azure Route Tables.",
			Attributes:  []string{"`table_name` (String, Required) Table identifier.", "`vpc_id` (String, Optional) Parent VPC ID."},
			BasicExample: `resource "multicloud_route_table" "basic" {
  provider_type = "aws"
  table_name    = "private-route-table"
}`,
			AdvancedExample: `resource "multicloud_route_table" "advanced" {
  provider_type = "azure"
  table_name    = "spoke-vnet-routes"

  extra_config = {
    "azure_disable_bgp_route_propagation" = "true"
  }
}`,
		},
		{
			Name:        "load_balancer",
			Category:    "Network",
			Description: "Unified Network / Application Load Balancer.",
			AWS:         "aws_lb",
			GCP:         "google_compute_forwarding_rule",
			Azure:       "azurerm_lb",
			HowItWorks:  "The `multicloud_load_balancer` resource provisions layer 4 (Network) and layer 7 (Application) load balancers across AWS ALB/NLB, GCP Forwarding Rules, and Azure Load Balancers.",
			Attributes:  []string{"`balancer_name` (String, Required) Load balancer name.", "`balancer_type` (String, Optional) 'application' or 'network'.", "`subnet_ids` (List[String], Optional) Target Subnet IDs."},
			BasicExample: `resource "multicloud_load_balancer" "basic" {
  provider_type = "gcp"
  balancer_name = "ingress-lb"
  balancer_type = "application"
}`,
			AdvancedExample: `resource "multicloud_load_balancer" "aws_advanced" {
  provider_type = "aws"
  balancer_name = "alb-external-aws"
  balancer_type = "application"

  extra_config = {
    "aws_idle_timeout"             = "60"
    "aws_enable_deletion_protection" = "true"
  }
}`,
		},
		{
			Name:        "api_gateway",
			Category:    "Network",
			Description: "Unified API Gateway Endpoint Manager.",
			AWS:         "aws_apigatewayv2_api",
			GCP:         "google_api_gateway_gateway",
			Azure:       "azurerm_api_management",
			HowItWorks:  "The `multicloud_api_gateway` resource provisions API gateway endpoints across AWS API Gateway v2, GCP API Gateway, and Azure API Management (APIM), exposing public HTTPS endpoints.",
			Attributes:  []string{"`api_name` (String, Required) API Gateway name.", "`protocol_type` (String, Optional) 'HTTP', 'REST', or 'WEBSOCKET'.", "`api_endpoint` (String, Read-Only) Computed HTTPS URL."},
			BasicExample: `resource "multicloud_api_gateway" "basic" {
  provider_type = "aws"
  api_name      = "customer-api"
  protocol_type = "HTTP"
}`,
			AdvancedExample: `resource "multicloud_api_gateway" "azure_advanced" {
  provider_type = "azure"
  api_name      = "apim-gateway-azure"
  protocol_type = "REST"

  extra_config = {
    "azure_publisher_email" = "admin@example.com"
    "azure_publisher_name"  = "Platform Team"
  }
}`,
		},
		{
			Name:        "cdn_distribution",
			Category:    "Network",
			Description: "Unified Content Delivery Network (CDN) Distribution.",
			AWS:         "aws_cloudfront_distribution",
			GCP:         "google_compute_backend_service",
			Azure:       "azurerm_cdn_endpoint",
			HowItWorks:  "The `multicloud_cdn_distribution` resource provisions edge caching content delivery distributions across AWS CloudFront, GCP Cloud CDN, and Azure CDN Endpoints.",
			Attributes:  []string{"`distribution_name` (String, Required) CDN name.", "`origin_domain` (String, Optional) Origin host domain.", "`enabled` (Bool, Optional) CDN status."},
			BasicExample: `resource "multicloud_cdn_distribution" "basic" {
  provider_type     = "aws"
  distribution_name = "assets-cdn"
  origin_domain     = "my-bucket.s3.amazonaws.com"
  enabled           = true
}`,
			AdvancedExample: `resource "multicloud_cdn_distribution" "advanced" {
  provider_type     = "aws"
  distribution_name = "global-media-cdn"
  origin_domain     = "media.example.com"
  enabled           = true

  extra_config = {
    "aws_price_class" = "PriceClass_100"
    "aws_ipv6_enabled" = "true"
  }
}`,
		},
		{
			Name:        "vpn_gateway",
			Category:    "Network",
			Description: "Unified Virtual Private Network (VPN) Gateway.",
			AWS:         "aws_vpn_gateway",
			GCP:         "google_compute_vpn_gateway",
			Azure:       "azurerm_virtual_network_gateway",
			HowItWorks:  "The `multicloud_vpn_gateway` resource provisions encrypted IPsec VPN gateways across AWS VPN Gateways, GCP Cloud VPN, and Azure VNet Gateways to connect hybrid on-premises sites to cloud networks.",
			Attributes:  []string{"`gateway_name` (String, Required) Gateway name.", "`vpc_id` (String, Optional) Target VPC ID.", "`tunnel_ip` (String, Optional) Remote tunnel IP."},
			BasicExample: `resource "multicloud_vpn_gateway" "basic" {
  provider_type = "aws"
  gateway_name  = "office-vpn"
}`,
			AdvancedExample: `resource "multicloud_vpn_gateway" "advanced" {
  provider_type = "gcp"
  gateway_name  = "hybrid-cloud-vpn"
  tunnel_ip     = "203.0.113.50"

  extra_config = {
    "gcp_ike_version" = "2"
  }
}`,
		},
		{
			Name:        "db_instance",
			Category:    "Database",
			Description: "Unified Relational Database Instance (SQL).",
			AWS:         "aws_db_instance",
			GCP:         "google_sql_database_instance",
			Azure:       "azurerm_postgresql_server",
			HowItWorks:  "The `multicloud_db_instance` resource provisions managed relational database instances across AWS RDS, GCP Cloud SQL, and Azure PostgreSQL/MySQL servers, configuring storage, engine versions, and Multi-AZ replication.",
			Attributes:  []string{"`instance_name` (String, Required) DB instance name.", "`engine` (String, Optional) 'postgres', 'mysql', or 'sqlserver'.", "`engine_version` (String, Optional) Version.", "`size_tier` (String, Optional) Size tier ('small', 'medium', 'large').", "`storage_gb` (Int64, Optional) Storage in GB.", "`multi_az` (Bool, Optional) High availability.", "`backup_retention_days` (Int64, Optional) Retention days."},
			BasicExample: `resource "multicloud_db_instance" "basic" {
  provider_type        = "aws"
  instance_name        = "prod-db"
  engine               = "postgres"
  engine_version       = "15.2"
  size_tier            = "medium"
  storage_gb           = 100
  multi_az             = true
  backup_retention_days = 7
}`,
			AdvancedExample: `resource "multicloud_db_instance" "gcp_advanced" {
  provider_type        = "gcp"
  instance_name        = "cloudsql-postgres-gcp"
  engine               = "postgres"
  engine_version       = "14"
  size_tier            = "large"
  storage_gb           = 250
  multi_az             = true
  backup_retention_days = 14

  extra_config = {
    "gcp_database_flags" = "max_connections=200"
  }
}`,
		},
		{
			Name:        "nosql_table",
			Category:    "Database",
			Description: "Unified NoSQL Document / Key-Value Database Table.",
			AWS:         "aws_dynamodb_table",
			GCP:         "google_firestore_database",
			Azure:       "azurerm_cosmosdb_account",
			HowItWorks:  "The `multicloud_nosql_table` resource provisions schemaless NoSQL tables across AWS DynamoDB, GCP Firestore, and Azure Cosmos DB.",
			Attributes:  []string{"`table_name` (String, Required) Table name.", "`hash_key` (String, Optional) Partition hash key.", "`range_key` (String, Optional) Sort range key.", "`billing_mode` (String, Optional) 'PAY_PER_REQUEST' or 'PROVISIONED'."},
			BasicExample: `resource "multicloud_nosql_table" "basic" {
  provider_type = "aws"
  table_name    = "user-sessions"
  hash_key      = "session_id"
  billing_mode  = "PAY_PER_REQUEST"
}`,
			AdvancedExample: `resource "multicloud_nosql_table" "azure_advanced" {
  provider_type = "azure"
  table_name    = "cosmos-events-azure"
  hash_key      = "event_id"

  extra_config = {
    "azure_consistency_level" = "Session"
  }
}`,
		},
		{
			Name:        "cache_cluster",
			Category:    "Database",
			Description: "Unified In-Memory Redis Cache Cluster.",
			AWS:         "aws_elasticache_cluster",
			GCP:         "google_redis_instance",
			Azure:       "azurerm_redis_cache",
			HowItWorks:  "The `multicloud_cache_cluster` resource provisions managed in-memory caching clusters across AWS ElastiCache, GCP Memorystore, and Azure Cache for Redis.",
			Attributes:  []string{"`cluster_name` (String, Required) Cache cluster name.", "`engine` (String, Optional) 'redis' or 'memcached'.", "`node_type` (String, Optional) Node tier.", "`num_nodes` (Int64, Optional) Node count."},
			BasicExample: `resource "multicloud_cache_cluster" "basic" {
  provider_type = "gcp"
  cluster_name  = "session-cache"
  engine        = "redis"
  num_nodes     = 2
}`,
			AdvancedExample: `resource "multicloud_cache_cluster" "aws_advanced" {
  provider_type = "aws"
  cluster_name  = "elasticache-redis-aws"
  engine        = "redis"
  node_type     = "cache.t4g.medium"
  num_nodes     = 3

  extra_config = {
    "aws_automatic_failover_enabled" = "true"
  }
}`,
		},
		{
			Name:        "data_warehouse",
			Category:    "Analytics",
			Description: "Unified Data Warehouse Analytical Cluster.",
			AWS:         "aws_redshift_cluster",
			GCP:         "google_bigquery_dataset",
			Azure:       "azurerm_synapse_workspace",
			HowItWorks:  "The `multicloud_data_warehouse` resource provisions petabyte-scale analytical data warehouse clusters across AWS Redshift, GCP BigQuery Datasets, and Azure Synapse Workspaces.",
			Attributes:  []string{"`warehouse_name` (String, Required) Analytical cluster name.", "`node_type` (String, Optional) Node spec.", "`num_nodes` (Int64, Optional) Node count.", "`database_name` (String, Optional) Initial database name."},
			BasicExample: `resource "multicloud_data_warehouse" "basic" {
  provider_type  = "aws"
  warehouse_name = "analytics-dw"
  node_type      = "ra3.4xlarge"
  num_nodes      = 4
}`,
			AdvancedExample: `resource "multicloud_data_warehouse" "gcp_advanced" {
  provider_type  = "gcp"
  warehouse_name = "bigquery_analytics_ds"

  extra_config = {
    "gcp_default_table_expiration_ms" = "3600000"
  }
}`,
		},
		{
			Name:        "search_index",
			Category:    "Analytics",
			Description: "Unified Search Engine Index Domain.",
			AWS:         "aws_opensearch_domain",
			GCP:         "google_discovery_engine_search_engine",
			Azure:       "azurerm_search_service",
			HowItWorks:  "The `multicloud_search_index` resource provisions managed search engine clusters across AWS OpenSearch, GCP Vertex AI Search, and Azure AI Search Services.",
			Attributes:  []string{"`index_name` (String, Required) Search index domain name.", "`instance_type` (String, Optional) Instance spec.", "`instance_count` (Int64, Optional) Node count."},
			BasicExample: `resource "multicloud_search_index" "basic" {
  provider_type  = "aws"
  index_name     = "app-logs-index"
  instance_count = 3
}`,
			AdvancedExample: `resource "multicloud_search_index" "advanced" {
  provider_type  = "aws"
  index_name     = "es-prod-cluster"
  instance_type  = "t3.medium.search"
  instance_count = 5

  extra_config = {
    "aws_dedicated_master_enabled" = "true"
  }
}`,
		},
		{
			Name:        "kubernetes_cluster",
			Category:    "Container",
			Description: "Unified Managed Kubernetes Engine (EKS, GKE, AKS).",
			AWS:         "aws_eks_cluster",
			GCP:         "google_container_cluster",
			Azure:       "azurerm_kubernetes_cluster",
			HowItWorks:  "The `multicloud_kubernetes_cluster` resource provisions managed Kubernetes control planes and node pools across AWS EKS, GCP GKE, and Azure AKS.",
			Attributes:  []string{"`cluster_name` (String, Required) Cluster name.", "`kubernetes_version` (String, Optional) K8s version.", "`node_count` (Int64, Optional) Worker node count.", "`node_instance_type` (String, Optional) Node tier."},
			BasicExample: `resource "multicloud_kubernetes_cluster" "basic" {
  provider_type      = "aws"
  cluster_name       = "prod-eks-cluster"
  kubernetes_version = "1.30"
  node_count         = 6
}`,
			AdvancedExample: `resource "multicloud_kubernetes_cluster" "gcp_advanced" {
  provider_type      = "gcp"
  cluster_name       = "gke-autopilot-cluster"
  kubernetes_version = "1.29"

  extra_config = {
    "gcp_enable_autopilot" = "true"
  }
}`,
		},
		{
			Name:        "container_registry",
			Category:    "Container",
			Description: "Unified Private Container Image Registry.",
			AWS:         "aws_ecr_repository",
			GCP:         "google_artifact_registry_repository",
			Azure:       "azurerm_container_registry",
			HowItWorks:  "The `multicloud_container_registry` resource provisions private OCI container registries across AWS ECR, GCP Artifact Registry, and Azure Container Registry (ACR).",
			Attributes:  []string{"`registry_name` (String, Required) Registry name.", "`image_tag_mutability` (String, Optional) 'MUTABLE' or 'IMMUTABLE'.", "`scan_on_push` (Bool, Optional) Scan on image push."},
			BasicExample: `resource "multicloud_container_registry" "basic" {
  provider_type        = "aws"
  registry_name        = "app-images"
  image_tag_mutability = "IMMUTABLE"
  scan_on_push         = true
}`,
			AdvancedExample: `resource "multicloud_container_registry" "azure_advanced" {
  provider_type        = "azure"
  registry_name        = "acrprodregistry"
  image_tag_mutability = "IMMUTABLE"

  extra_config = {
    "azure_admin_enabled" = "false"
    "azure_sku"           = "Premium"
  }
}`,
		},
		{
			Name:        "security_group",
			Category:    "Security",
			Description: "Unified Firewall Security Group.",
			AWS:         "aws_security_group",
			GCP:         "google_compute_firewall",
			Azure:       "azurerm_network_security_group",
			HowItWorks:  "The `multicloud_security_group` resource provisions virtual firewall rules across AWS Security Groups, GCP Compute Firewalls, and Azure Network Security Groups (NSGs).",
			Attributes:  []string{"`group_name` (String, Required) Security group name.", "`description` (String, Optional) Group description.", "`vpc_id` (String, Optional) Parent VPC ID."},
			BasicExample: `resource "multicloud_security_group" "basic" {
  provider_type = "aws"
  group_name    = "web-sec-group"
  description   = "Allow HTTPS inbound traffic"
}`,
			AdvancedExample: `resource "multicloud_security_group" "gcp_advanced" {
  provider_type = "gcp"
  group_name    = "allow-internal-traffic"
  description   = "GCP Compute internal firewall"

  extra_config = {
    "gcp_target_tags" = "web-node,app-node"
  }
}`,
		},
		{
			Name:        "secret",
			Category:    "Security",
			Description: "Unified Key-Value Secret Vault Store.",
			AWS:         "aws_secretsmanager_secret",
			GCP:         "google_secret_manager_secret",
			Azure:       "azurerm_key_vault_secret",
			HowItWorks:  "The `multicloud_secret` resource provisions encrypted key-value secrets in AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault.",
			Attributes:  []string{"`secret_name` (String, Required) Secret name.", "`description` (String, Optional) Secret description.", "`secret_string` (String, Optional, Sensitive) Secret payload."},
			BasicExample: `resource "multicloud_secret" "basic" {
  provider_type = "gcp"
  secret_name   = "db-password"
  secret_string = var.db_password
}`,
			AdvancedExample: `resource "multicloud_secret" "aws_advanced" {
  provider_type = "aws"
  secret_name   = "third-party-api-key"
  secret_string = var.api_key

  extra_config = {
    "aws_recovery_window_in_days" = "7"
  }
}`,
		},
		{
			Name:        "kms_key",
			Category:    "Security",
			Description: "Unified KMS Encryption Key.",
			AWS:         "aws_kms_key",
			GCP:         "google_kms_crypto_key",
			Azure:       "azurerm_key_vault_key",
			HowItWorks:  "The `multicloud_kms_key` resource provisions Customer-Managed KMS Encryption Keys (CMKs) across AWS KMS, GCP KMS, and Azure Key Vault, managing deletion windows and automated key rotation.",
			Attributes:  []string{"`key_name` (String, Required) Key alias name.", "`description` (String, Optional) Key description.", "`deletion_window_in_days` (Int64, Optional) Deletion window.", "`enable_key_rotation` (Bool, Optional) Enable key rotation."},
			BasicExample: `resource "multicloud_kms_key" "basic" {
  provider_type           = "aws"
  key_name                = "app-kms-key"
  deletion_window_in_days = 30
  enable_key_rotation     = true
}`,
			AdvancedExample: `resource "multicloud_kms_key" "gcp_advanced" {
  provider_type       = "gcp"
  key_name            = "gcp-cmk-key"
  enable_key_rotation = true

  extra_config = {
    "gcp_rotation_period" = "7776000s"
  }
}`,
		},
		{
			Name:        "iam_role",
			Category:    "IAM",
			Description: "Unified IAM Role and Identity Policy.",
			AWS:         "aws_iam_role",
			GCP:         "google_service_account",
			Azure:       "azurerm_user_assigned_identity",
			HowItWorks:  "The `multicloud_iam_role` resource provisions identity roles across AWS IAM Roles, GCP Service Accounts, and Azure Managed Identities.",
			Attributes:  []string{"`role_name` (String, Required) Role name.", "`assume_role_policy` (String, Optional) Trust policy JSON.", "`description` (String, Optional) Role description."},
			BasicExample: `resource "multicloud_iam_role" "basic" {
  provider_type = "aws"
  role_name     = "app-execution-role"
}`,
			AdvancedExample: `resource "multicloud_iam_role" "aws_advanced" {
  provider_type       = "aws"
  role_name           = "ecs-task-execution-role"
  assume_role_policy  = file("${path.module}/policies/trust_policy.json")

  extra_config = {
    "aws_max_session_duration" = "7200"
  }
}`,
		},
		{
			Name:        "dns_zone",
			Category:    "DNS",
			Description: "Unified Managed DNS Zone.",
			AWS:         "aws_route53_zone",
			GCP:         "google_dns_managed_zone",
			Azure:       "azurerm_dns_zone",
			HowItWorks:  "The `multicloud_dns_zone` resource provisions authoritative DNS zones across AWS Route 53, GCP Cloud DNS, and Azure DNS.",
			Attributes:  []string{"`zone_name` (String, Required) DNS zone name.", "`domain_name` (String, Optional) FQDN domain name."},
			BasicExample: `resource "multicloud_dns_zone" "basic" {
  provider_type = "aws"
  zone_name     = "prod-zone"
  domain_name   = "example.com"
}`,
			AdvancedExample: `resource "multicloud_dns_zone" "gcp_advanced" {
  provider_type = "gcp"
  zone_name     = "corp-internal-dns"
  domain_name   = "internal.company.net"

  extra_config = {
    "gcp_visibility" = "private"
  }
}`,
		},
		{
			Name:        "pubsub_topic",
			Category:    "Messaging",
			Description: "Unified Publish/Subscribe Topic.",
			AWS:         "aws_sns_topic",
			GCP:         "google_pubsub_topic",
			Azure:       "azurerm_servicebus_topic",
			HowItWorks:  "The `multicloud_pubsub_topic` resource provisions pub/sub messaging topics across AWS SNS, GCP Pub/Sub, and Azure Service Bus Topics.",
			Attributes:  []string{"`topic_name` (String, Required) Topic name.", "`display_name` (String, Optional) Friendly label."},
			BasicExample: `resource "multicloud_pubsub_topic" "basic" {
  provider_type = "gcp"
  topic_name    = "user-signup-events"
}`,
			AdvancedExample: `resource "multicloud_pubsub_topic" "aws_advanced" {
  provider_type = "aws"
  topic_name    = "payment-alerts-sns"
  display_name  = "Payment Processing Alerts"

  extra_config = {
    "aws_kms_master_key_id" = "alias/aws/sns"
  }
}`,
		},
		{
			Name:        "message_queue",
			Category:    "Messaging",
			Description: "Unified Message Queue.",
			AWS:         "aws_sqs_queue",
			GCP:         "google_pubsub_subscription",
			Azure:       "azurerm_servicebus_queue",
			HowItWorks:  "The `multicloud_message_queue` resource provisions asynchronous message queues across AWS SQS, GCP Pub/Sub Subscriptions, and Azure Service Bus Queues.",
			Attributes:  []string{"`queue_name` (String, Required) Queue name.", "`delay_seconds` (Int64, Optional) Delivery delay.", "`max_message_size` (Int64, Optional) Max size in bytes.", "`message_retention_seconds` (Int64, Optional) Retention seconds."},
			BasicExample: `resource "multicloud_message_queue" "basic" {
  provider_type             = "aws"
  queue_name                = "order-processing-queue"
  delay_seconds             = 0
  message_retention_seconds = 86400
}`,
			AdvancedExample: `resource "multicloud_message_queue" "aws_fifo_advanced" {
  provider_type             = "aws"
  queue_name                = "financial-transactions.fifo"
  message_retention_seconds = 1209600

  extra_config = {
    "aws_fifo_queue"                  = "true"
    "aws_content_based_deduplication" = "true"
  }
}`,
		},
		{
			Name:        "event_bridge",
			Category:    "Messaging",
			Description: "Unified Event Router & Trigger Bus.",
			AWS:         "aws_cloudwatch_event_bus",
			GCP:         "google_eventarc_trigger",
			Azure:       "azurerm_eventgrid_system_topic",
			HowItWorks:  "The `multicloud_event_bridge` resource provisions event routing buses across AWS EventBridge, GCP Eventarc, and Azure Event Grid System Topics.",
			Attributes:  []string{"`bus_name` (String, Required) Event bus router name.", "`event_source` (String, Optional) Event source identifier."},
			BasicExample: `resource "multicloud_event_bridge" "basic" {
  provider_type = "aws"
  bus_name      = "custom-app-eventbus"
}`,
			AdvancedExample: `resource "multicloud_event_bridge" "gcp_advanced" {
  provider_type = "gcp"
  bus_name      = "cloud-storage-trigger"

  extra_config = {
    "gcp_destination_run_service" = "event-processor"
  }
}`,
		},
		{
			Name:        "monitoring_dashboard",
			Category:    "Observability",
			Description: "Unified Cloud Monitoring Dashboard.",
			AWS:         "aws_cloudwatch_dashboard",
			GCP:         "google_monitoring_dashboard",
			Azure:       "azurerm_portal_dashboard",
			HowItWorks:  "The `multicloud_monitoring_dashboard` resource provisions visual metrics dashboards across AWS CloudWatch, GCP Cloud Monitoring, and Azure Portal Dashboards.",
			Attributes:  []string{"`dashboard_name` (String, Required) Dashboard name.", "`dashboard_body` (String, Optional) JSON widget layout."},
			BasicExample: `resource "multicloud_monitoring_dashboard" "basic" {
  provider_type  = "aws"
  dashboard_name = "platform-overview"
}`,
			AdvancedExample: `resource "multicloud_monitoring_dashboard" "azure_advanced" {
  provider_type  = "azure"
  dashboard_name = "k8s-metrics-dashboard"

  extra_config = {
    "azure_dashboard_type" = "shared"
  }
}`,
		},
	}

	for _, r := range resources {
		attrDocs := ""
		for _, attr := range r.Attributes {
			attrDocs += fmt.Sprintf("- %s\n", attr)
		}
		attrDocs += "- `region` (String, Optional) Target placement region.\n"
		attrDocs += "- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.\n"

		resourceDoc := fmt.Sprintf(`---
# subcategory: "%s"
page_title: "multicloud_%s Resource - terraform-provider-multicloud"
description: |-
  %s
---

# multicloud_%s (Resource)

%s

## Cloud Targets
- **AWS Target:** %s
- **GCP Target:** %s
- **Azure Target:** %s

## How It Works

%s

## Example Usage

### Basic Usage
`+"```hcl"+`
%s
`+"```"+`

### Advanced Usage with Cloud Escape Hatches (`+"`extra_config`"+`)
`+"```hcl"+`
%s
`+"```"+`

## Schema Attributes

### Required
- `+"`provider_type`"+` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
%s
### Read-Only
- `+"`id`"+` (String) State resource identifier (<cloud>/<region>/<name>).
`, r.Category, r.Name, r.Description, r.Name, r.Description, r.AWS, r.GCP, r.Azure, r.HowItWorks, r.BasicExample, r.AdvancedExample, attrDocs)

		filename := fmt.Sprintf("%s.md", r.Name)
		docPath := filepath.Join(docsDir, filename)
		/* #nosec G306 G703 */
		if err := os.WriteFile(docPath, []byte(resourceDoc), 0600); err != nil {
			fmt.Printf("Error writing %s: %v\n", filename, err)
			os.Exit(1)
		}
	}

	fmt.Printf("[Doc Generator] Successfully generated multi-example registry documentation with 'How It Works' sections for all %d unified resources in %s/\n", len(resources), docsDir)
}
