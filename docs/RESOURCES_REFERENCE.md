# Multi-Cloud Terraform Provider (`terraform-provider-multicloud`) Resources Reference Manual

This technical manual details the complete schema attributes, required/optional parameters, read-only values, cloud targets, and `extra_config` escape hatches for all **43 unified resources** provided by `terraform-provider-multicloud`.

---

## Common Attributes Across All Resources

All unified `multicloud_*` resources share the following standard attributes:

- `provider_type` (String, Required) - Target cloud provider (`'aws'`, `'gcp'`, or `'azure'`).
- `region` (String, Optional) - Target placement region (e.g. `'us-west-2'`, `'us-central1'`, `'eastus'`).
- `id` (String, Read-Only) - State resource identifier string (`<cloud>/<region>/<name>`).
- `extra_config` (Map[String], Optional) - Cloud-specific escape-hatch key-value map passed through to upstream cloud SDKs.

---

## 1. Storage & Replication

### `multicloud_storage_bucket`
Unified object storage container.
- **AWS Target:** `aws_s3_bucket`
- **GCP Target:** `google_storage_bucket`
- **Azure Target:** `azurerm_storage_container`
- **Attributes:**
  - `bucket_name` (String, Required) - Bucket identifier.
  - `versioning_enabled` (Bool, Optional) - Enable object versioning.
  - `encryption_enabled` (Bool, Optional) - Enable server-side encryption.
  - `extra_config` (Map[String], Optional) - Escape hatches (e.g. `"aws_force_destroy" = "true"`, `"gcp_storage_class" = "NEARLINE"`).

### `multicloud_data_sync`
Multi-cloud object storage background data synchronization.
- **AWS Target:** S3 Cross-Region Replication
- **GCP Target:** Storage Transfer Service
- **Azure Target:** Azure Storage Sync Service
- **Attributes:**
  - `sync_name` (String, Required) - Data sync task name.
  - `source_bucket` (String, Optional) - Source bucket URI.
  - `destination_bucket` (String, Optional) - Destination bucket URI.
  - `schedule_cron` (String, Optional) - Cron schedule expression.

---

## 2. Compute & Scaling

### `multicloud_virtual_machine`
Unified virtual machine compute instance.
- **AWS Target:** `aws_instance`
- **GCP Target:** `google_compute_instance`
- **Azure Target:** `azurerm_linux_virtual_machine`
- **Attributes:**
  - `vm_name` (String, Required) - Instance name.
  - `size_tier` (String, Optional) - Instance tier (`'small'`, `'medium'`, `'large'`).
  - `image_id` (String, Optional) - Custom image/AMI ID.
  - `subnet_id` (String, Optional) - VPC Subnet ID.
  - `ssh_public_key` (String, Optional, Sensitive) - SSH public key.
  - `tags` (Map[String], Optional) - Key-value tags.

### `multicloud_auto_scaling_group`
Unified compute auto-scaling group.
- **AWS Target:** `aws_autoscaling_group`
- **GCP Target:** `google_compute_instance_group_manager`
- **Azure Target:** `azurerm_linux_virtual_machine_scale_set`
- **Attributes:**
  - `group_name` (String, Required) - ASG group name.
  - `min_size` (Int64, Optional) - Minimum instance capacity.
  - `max_size` (Int64, Optional) - Maximum instance capacity.
  - `desired_capacity` (Int64, Optional) - Desired target capacity.

### `multicloud_serverless_function`
Unified serverless event function.
- **AWS Target:** `aws_lambda_function`
- **GCP Target:** `google_cloudfunctions2_function`
- **Azure Target:** `azurerm_linux_function_app`
- **Attributes:**
  - `function_name` (String, Required) - Function name.
  - `runtime` (String, Optional) - Language runtime (`'nodejs20'`, `'python3.11'`).
  - `handler` (String, Optional) - Entrypoint handler function.
  - `memory_size_mb` (Int64, Optional) - Allocated memory in MB.
  - `timeout_seconds` (Int64, Optional) - Execution timeout in seconds.
  - `environment_variables` (Map[String], Optional) - Runtime environment variables.

---

## 3. Network & Infrastructure

### `multicloud_virtual_network`
Unified virtual private network.
- **AWS Target:** `aws_vpc`
- **GCP Target:** `google_compute_network`
- **Azure Target:** `azurerm_virtual_network`
- **Attributes:**
  - `network_name` (String, Required) - Network name.
  - `cidr_block` (String, Optional) - Primary IPv4 CIDR range (e.g. `'10.0.0.0/16'`).

### `multicloud_subnet`
Unified subnetwork.
- **AWS Target:** `aws_subnet`
- **GCP Target:** `google_compute_subnetwork`
- **Azure Target:** `azurerm_subnet`
- **Attributes:**
  - `subnet_name` (String, Required) - Subnet name.
  - `vpc_id` (String, Optional) - Parent Virtual Network ID.
  - `cidr_block` (String, Optional) - Subnet IPv4 CIDR range.
  - `availability_zone` (String, Optional) - AZ placement.

### `multicloud_static_ip`
Unified static elastic IP address.
- **AWS Target:** `aws_eip`
- **GCP Target:** `google_compute_address`
- **Azure Target:** `azurerm_public_ip`
- **Attributes:**
  - `ip_name` (String, Required) - IP allocation name.
  - `allocation_type` (String, Optional) - `'static'` or `'dynamic'`.

### `multicloud_nat_gateway`
Unified NAT gateway.
- **AWS Target:** `aws_nat_gateway`
- **GCP Target:** `google_compute_router_nat`
- **Azure Target:** `azurerm_nat_gateway`
- **Attributes:**
  - `gateway_name` (String, Required) - NAT gateway name.
  - `subnet_id` (String, Optional) - Target Subnet ID.
  - `allocation_id` (String, Optional) - Static IP Allocation ID.

### `multicloud_route_table`
Unified routing table.
- **AWS Target:** `aws_route_table`
- **GCP Target:** `google_compute_route`
- **Azure Target:** `azurerm_route_table`
- **Attributes:**
  - `table_name` (String, Required) - Route table name.
  - `vpc_id` (String, Optional) - Parent VPC ID.

### `multicloud_load_balancer`
Unified network / application load balancer.
- **AWS Target:** `aws_lb`
- **GCP Target:** `google_compute_forwarding_rule`
- **Azure Target:** `azurerm_lb`
- **Attributes:**
  - `balancer_name` (String, Required) - Load balancer name.
  - `balancer_type` (String, Optional) - `'application'` or `'network'`.
  - `subnet_ids` (List[String], Optional) - Target subnet IDs.

### `multicloud_api_gateway`
Unified API Gateway endpoint manager.
- **AWS Target:** `aws_apigatewayv2_api`
- **GCP Target:** `google_api_gateway_gateway`
- **Azure Target:** `azurerm_api_management`
- **Attributes:**
  - `api_name` (String, Required) - API Gateway name.
  - `protocol_type` (String, Optional) - `'HTTP'`, `'REST'`, or `'WEBSOCKET'`.
  - `api_endpoint` (String, Read-Only) - Computed HTTPS API endpoint URL.

### `multicloud_cdn_distribution`
Unified Content Delivery Network (CDN) distribution.
- **AWS Target:** `aws_cloudfront_distribution`
- **GCP Target:** `google_compute_backend_service`
- **Azure Target:** `azurerm_cdn_endpoint`
- **Attributes:**
  - `distribution_name` (String, Required) - CDN Distribution name.
  - `origin_domain` (String, Optional) - Origin domain host.
  - `enabled` (Bool, Optional) - Distribution status.

### `multicloud_vpn_gateway`
Unified Virtual Private Network (VPN) gateway.
- **AWS Target:** `aws_vpn_gateway`
- **GCP Target:** `google_compute_vpn_gateway`
- **Azure Target:** `azurerm_virtual_network_gateway`
- **Attributes:**
  - `gateway_name` (String, Required) - VPN Gateway name.
  - `vpc_id` (String, Optional) - Target VPC ID.
  - `tunnel_ip` (String, Optional) - Remote Tunnel IP.

---

## 4. Databases & Storage

### `multicloud_db_instance`
Unified relational database instance (SQL).
- **AWS Target:** `aws_db_instance`
- **GCP Target:** `google_sql_database_instance`
- **Azure Target:** `azurerm_postgresql_server`
- **Attributes:**
  - `instance_name` (String, Required) - DB instance name.
  - `engine` (String, Optional) - `'postgres'`, `'mysql'`, or `'sqlserver'`.
  - `engine_version` (String, Optional) - Database version.
  - `size_tier` (String, Optional) - DB tier (`'small'`, `'medium'`, `'large'`).
  - `storage_gb` (Int64, Optional) - Storage size in GB.
  - `multi_az` (Bool, Optional) - High availability Multi-AZ deployment.
  - `backup_retention_days` (Int64, Optional) - Automated backup retention days.

### `multicloud_nosql_table`
Unified NoSQL document / key-value database table.
- **AWS Target:** `aws_dynamodb_table`
- **GCP Target:** `google_firestore_database`
- **Azure Target:** `azurerm_cosmosdb_account`
- **Attributes:**
  - `table_name` (String, Required) - Table name.
  - `hash_key` (String, Optional) - Partition hash key.
  - `range_key` (String, Optional) - Sort range key.
  - `billing_mode` (String, Optional) - `'PAY_PER_REQUEST'` or `'PROVISIONED'`.

### `multicloud_cache_cluster`
Unified in-memory Redis cache cluster.
- **AWS Target:** `aws_elasticache_cluster`
- **GCP Target:** `google_redis_instance`
- **Azure Target:** `azurerm_redis_cache`
- **Attributes:**
  - `cluster_name` (String, Required) - Cache cluster name.
  - `engine` (String, Optional) - `'redis'` or `'memcached'`.
  - `node_type` (String, Optional) - Cache node tier.
  - `num_nodes` (Int64, Optional) - Cache node count.

---

## 5. Analytics & Search

### `multicloud_data_warehouse`
Unified analytical data warehouse cluster.
- **AWS Target:** `aws_redshift_cluster`
- **GCP Target:** `google_bigquery_dataset`
- **Azure Target:** `azurerm_synapse_workspace`
- **Attributes:**
  - `warehouse_name` (String, Required) - Analytical cluster name.
  - `node_type` (String, Optional) - Compute node spec.
  - `num_nodes` (Int64, Optional) - Node count.
  - `database_name` (String, Optional) - Initial database name.

### `multicloud_search_index`
Unified search engine index domain.
- **AWS Target:** `aws_opensearch_domain`
- **GCP Target:** `google_discovery_engine_search_engine`
- **Azure Target:** `azurerm_search_service`
- **Attributes:**
  - `index_name` (String, Required) - Search index domain name.
  - `instance_type` (String, Optional) - Search node instance spec.
  - `instance_count` (Int64, Optional) - Node count.

---

## 6. Containers & IDP

### `multicloud_kubernetes_cluster`
Unified managed Kubernetes engine (EKS / GKE / AKS).
- **AWS Target:** `aws_eks_cluster`
- **GCP Target:** `google_container_cluster`
- **Azure Target:** `azurerm_kubernetes_cluster`
- **Attributes:**
  - `cluster_name` (String, Required) - Cluster name.
  - `kubernetes_version` (String, Optional) - K8s control plane version.
  - `node_count` (Int64, Optional) - Node pool size.
  - `node_instance_type` (String, Optional) - Worker node instance tier.

### `multicloud_container_registry`
Unified private container image registry.
- **AWS Target:** `aws_ecr_repository`
- **GCP Target:** `google_artifact_registry_repository`
- **Azure Target:** `azurerm_container_registry`
- **Attributes:**
  - `registry_name` (String, Required) - Registry repository name.
  - `image_tag_mutability` (String, Optional) - `'MUTABLE'` or `'IMMUTABLE'`.
  - `scan_on_push` (Bool, Optional) - Vulnerability scan on image push.

---

## 7. Security, IAM, DNS, Messaging & Reliability

### `multicloud_security_group`
Unified firewall security group.
- **AWS Target:** `aws_security_group`
- **GCP Target:** `google_compute_firewall`
- **Azure Target:** `azurerm_network_security_group`
- **Attributes:**
  - `group_name` (String, Required) - Security group name.
  - `description` (String, Optional) - Security group description.
  - `vpc_id` (String, Optional) - Parent VPC ID.

### `multicloud_secret`
Unified key-value secret vault store.
- **AWS Target:** `aws_secretsmanager_secret`
- **GCP Target:** `google_secret_manager_secret`
- **Azure Target:** `azurerm_key_vault_secret`
- **Attributes:**
  - `secret_name` (String, Required) - Secret vault name.
  - `description` (String, Optional) - Secret description.
  - `secret_string` (String, Optional, Sensitive) - Secret payload.

### `multicloud_kms_key`
Unified KMS encryption key.
- **AWS Target:** `aws_kms_key`
- **GCP Target:** `google_kms_crypto_key`
- **Azure Target:** `azurerm_key_vault_key`
- **Attributes:**
  - `key_name` (String, Required) - Key alias name.
  - `description` (String, Optional) - Key description.
  - `deletion_window_in_days` (Int64, Optional) - Deletion grace period.
  - `enable_key_rotation` (Bool, Optional) - Enable automated annual key rotation.

### `multicloud_iam_role`
Unified IAM role and identity policy.
- **AWS Target:** `aws_iam_role`
- **GCP Target:** `google_service_account`
- **Azure Target:** `azurerm_user_assigned_identity`
- **Attributes:**
  - `role_name` (String, Required) - Role name.
  - `assume_role_policy` (String, Optional) - Trust policy JSON document.
  - `description` (String, Optional) - Role description.

### `multicloud_dns_zone`
Unified managed DNS zone.
- **AWS Target:** `aws_route53_zone`
- **GCP Target:** `google_dns_managed_zone`
- **Azure Target:** `azurerm_dns_zone`
- **Attributes:**
  - `zone_name` (String, Required) - DNS zone identifier.
  - `domain_name` (String, Optional) - FQDN domain name.

### `multicloud_pubsub_topic`
Unified Publish/Subscribe topic.
- **AWS Target:** `aws_sns_topic`
- **GCP Target:** `google_pubsub_topic`
- **Azure Target:** `azurerm_servicebus_topic`
- **Attributes:**
  - `topic_name` (String, Required) - Topic name.
  - `display_name` (String, Optional) - Friendly display label.

### `multicloud_message_queue`
Unified message queue.
- **AWS Target:** `aws_sqs_queue`
- **GCP Target:** `google_pubsub_subscription`
- **Azure Target:** `azurerm_servicebus_queue`
- **Attributes:**
  - `queue_name` (String, Required) - Queue name.
  - `delay_seconds` (Int64, Optional) - Message delivery delay.
  - `max_message_size` (Int64, Optional) - Maximum message size in bytes.
  - `message_retention_seconds` (Int64, Optional) - Message retention window.

### `multicloud_event_bridge`
Unified event router and trigger bus.
- **AWS Target:** `aws_cloudwatch_event_bus`
- **GCP Target:** `google_eventarc_trigger`
- **Azure Target:** `azurerm_eventgrid_system_topic`
- **Attributes:**
  - `bus_name` (String, Required) - Event bus router name.
  - `event_source` (String, Optional) - Event source identifier.

### `multicloud_failover_policy`
Unified active-passive disaster recovery failover policy.
- **AWS Target:** Route53 Failover Routing
- **GCP Target:** Cloud DNS Failover Policy
- **Azure Target:** Azure Traffic Manager / Front Door
- **Attributes:**
  - `policy_name` (String, Required) - Failover policy name.
  - `primary_cloud` (String, Required) - Primary cloud target (`'aws'`, `'gcp'`, `'azure'`).
  - `failover_cloud` (String, Required) - Backup failover cloud target.
  - `health_check_url` (String, Optional) - Endpoint health check URL.
  - `auto_failover` (Bool, Optional) - Enable automated failover.
  - `failover_status` (String, Read-Only) - Computed status (`'PRIMARY_HEALTHY'`).

### `multicloud_identity_federation`
Unified OIDC workload identity federation.
- **AWS Target:** IAM OIDC Provider
- **GCP Target:** Workload Identity
- **Azure Target:** Entra ID Workload Identity
- **Attributes:**
  - `federation_name` (String, Required) - Identity federation name.
  - `issuer_url` (String, Optional) - External OIDC Issuer URL.
  - `client_id_list` (List[String], Optional) - Target audience client IDs.

### `multicloud_secret_rotator`
Unified automated secret key rotator.
- **AWS Target:** Secrets Manager Rotation
- **GCP Target:** Secret Manager Rotation
- **Azure Target:** Key Vault Rotation
- **Attributes:**
  - `rotator_name` (String, Required) - Rotation schedule name.
  - `secret_id` (String, Optional) - Target Secret Vault ID.
  - `rotation_days` (Int64, Optional) - Rotation interval in days.

### `multicloud_monitoring_dashboard`
Unified cloud monitoring dashboard.
- **AWS Target:** CloudWatch Dashboard
- **GCP Target:** Cloud Monitoring Dashboard
- **Azure Target:** Azure Portal Dashboard
- **Attributes:**
  - `dashboard_name` (String, Required) - Dashboard name.
  - `dashboard_body` (String, Optional) - JSON dashboard widget definition.

---

## 8. New Enterprise Extensions

### `multicloud_container_app`
Unified serverless container workload app.
- **AWS Target:** `aws_apprunner_service` / ECS Fargate
- **GCP Target:** `google_cloud_run_v2_service`
- **Azure Target:** `azurerm_container_app`
- **Attributes:**
  - `app_name` (String, Required) - App service name.
  - `image` (String, Optional) - Container image URI.
  - `cpu` (String, Optional) - CPU specification.
  - `memory` (String, Optional) - Memory specification.
  - `port` (Int64, Optional) - Container port.

### `multicloud_bastion_host`
Unified managed SSH bastion jump host.
- **AWS Target:** `aws_ec2_instance_connect_endpoint`
- **GCP Target:** `google_iap_tunnel`
- **Azure Target:** `azurerm_bastion_host`
- **Attributes:**
  - `host_name` (String, Required) - Bastion host name.
  - `vpc_id` (String, Optional) - Parent VPC ID.
  - `subnet_id` (String, Optional) - Target Subnet ID.

### `multicloud_waf_policy`
Unified Web Application Firewall policy.
- **AWS Target:** `aws_wafv2_web_acl`
- **GCP Target:** `google_compute_security_policy` (Cloud Armor)
- **Azure Target:** `azurerm_web_application_firewall_policy`
- **Attributes:**
  - `policy_name` (String, Required) - WAF policy name.
  - `default_action` (String, Optional) - Default action ('allow' or 'block').

### `multicloud_vpc_peering`
Unified Virtual Private Network Peering Connection.
- **AWS Target:** `aws_vpc_peering_connection`
- **GCP Target:** `google_compute_network_peering`
- **Azure Target:** `azurerm_virtual_network_peering`
- **Attributes:**
  - `peering_name` (String, Required) - Peering connection name.
  - `vpc_id` (String, Optional) - Local VPC ID.
  - `peer_vpc_id` (String, Optional) - Remote VPC ID.
  - `peer_region` (String, Optional) - Remote VPC region.

### `multicloud_app_config`
Unified key-value application configuration store.
- **AWS Target:** `aws_ssm_parameter`
- **GCP Target:** `google_runtimeconfig_config`
- **Azure Target:** `azurerm_app_configuration`
- **Attributes:**
  - `config_name` (String, Required) - Config store name.
  - `config_key` (String, Optional) - Parameter key.
  - `config_value` (String, Optional) - Parameter value.

### `multicloud_ai_endpoint`
Unified AI / Machine Learning inference model endpoint.
- **AWS Target:** `aws_sagemaker_endpoint` / Bedrock
- **GCP Target:** `google_vertex_ai_endpoint`
- **Azure Target:** `azurerm_cognitive_account` (Azure OpenAI)
- **Attributes:**
  - `endpoint_name` (String, Required) - Inference endpoint name.
  - `model_name` (String, Optional) - Deployed model identifier.
  - `instance_type` (String, Optional) - Compute hardware spec.

### `multicloud_streaming_cluster`
Unified managed Apache Kafka event streaming cluster.
- **AWS Target:** `aws_msk_cluster`
- **GCP Target:** `google_managed_kafka_cluster`
- **Azure Target:** `azurerm_eventhub_namespace`
- **Attributes:**
  - `cluster_name` (String, Required) - Streaming cluster name.
  - `kafka_version` (String, Optional) - Kafka version.
  - `node_count` (Int64, Optional) - Broker node count.

### `multicloud_metric_alert`
Unified metric threshold alarm rule.
- **AWS Target:** `aws_cloudwatch_metric_alarm`
- **GCP Target:** `google_monitoring_alert_policy`
- **Azure Target:** `azurerm_monitor_metric_alert`
- **Attributes:**
  - `alert_name` (String, Required) - Alert rule name.
  - `metric_name` (String, Optional) - Target metric.
  - `threshold` (Float64, Optional) - Evaluation threshold.
  - `comparison` (String, Optional) - Operator.

### `multicloud_log_workspace`
Unified centralized log analytics workspace.
- **AWS Target:** `aws_cloudwatch_log_group`
- **GCP Target:** `google_logging_project_sink`
- **Azure Target:** `azurerm_log_analytics_workspace`
- **Attributes:**
  - `workspace_name` (String, Required) - Log workspace name.
  - `retention_days` (Int64, Optional) - Log retention window in days.

### `multicloud_graphql_api`
Unified managed GraphQL API endpoint.
- **AWS Target:** `aws_appsync_graphql_api`
- **GCP Target:** `google_apigee_environment`
- **Azure Target:** `azurerm_api_management_api`
- **Attributes:**
  - `api_name` (String, Required) - GraphQL API name.
  - `authentication_type` (String, Optional) - Auth type.
  - `schema_definition` (String, Optional) - GraphQL SDL schema string.

