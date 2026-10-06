# Multi-Cloud Terraform Provider (`terraform-provider-multicloud`) Resources Reference Manual

This technical manual details the complete schema attributes, required/optional parameters, read-only values, cloud targets, and `extra_config` escape hatches for all **70 unified resources** and **2 data sources** provided by `terraform-provider-multicloud`.

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

### `multicloud_dns_record`
Unified DNS resource record set.
- **AWS Target:** `aws_route53_record`
- **GCP Target:** `google_dns_record_set`
- **Azure Target:** `azurerm_dns_a_record`
- **Attributes:**
  - `zone_id` (String, Required) - Parent DNS zone ID.
  - `record_name` (String, Required) - Record hostname.
  - `record_type` (String, Required) - DNS type (A, CNAME, etc.).
  - `ttl` (Int64, Optional) - Time-To-Live.
  - `records` (List[String], Optional) - List of target IPs or hostnames.

### `multicloud_dns_health_check`
Unified automated DNS health check monitor.
- **AWS Target:** `aws_route53_health_check`
- **GCP Target:** `google_monitoring_uptime_check_config`
- **Azure Target:** `azurerm_traffic_manager_endpoint`
- **Attributes:**
  - `check_name` (String, Required) - Health check probe name.
  - `type` (String, Required) - Probe protocol (HTTP, HTTPS, TCP).
  - `fqdn` (String, Optional) - Target FQDN.
  - `ip_address` (String, Optional) - Target IP address.
  - `port` (Int64, Optional) - Target port.
  - `resource_path` (String, Optional) - HTTP URI path.
  - `failure_threshold` (Int64, Optional) - Threshold count.

### `multicloud_dns_zone_link`
Unified private DNS network association.
- **AWS Target:** `aws_route53_zone_association`
- **GCP Target:** `google_dns_managed_zone` (private visibility)
- **Azure Target:** `azurerm_private_dns_zone_virtual_network_link`
- **Attributes:**
  - `link_name` (String, Required) - Association name.
  - `zone_id` (String, Required) - Private DNS zone ID.
  - `vpc_id` (String, Required) - Target VPC/VNet ID.
  - `registration_enabled` (Bool, Optional) - Auto-register VM hostnames.

### `multicloud_dns_resolver`
Unified hybrid cloud DNS resolver endpoint.
- **AWS Target:** `aws_route53_resolver_endpoint`
- **GCP Target:** `google_dns_policy`
- **Azure Target:** `azurerm_private_dns_resolver`
- **Attributes:**
  - `resolver_name` (String, Required) - Resolver endpoint name.
  - `direction` (String, Required) - INBOUND or OUTBOUND.
  - `vpc_id` (String, Required) - Parent VPC ID.
  - `ip_configurations` (List[String], Optional) - Subnet/IP list.

### `multicloud_dnssec`
Unified DNSSEC zone signing security.
- **AWS Target:** `aws_route53_key_signing_key`
- **GCP Target:** `google_dns_managed_zone`
- **Azure Target:** `azurerm_dns_zone`
- **Attributes:**
  - `zone_id` (String, Required) - Target DNS zone ID.
  - `state` (String, Required) - Signing state (ON/OFF).
  - `key_type` (String, Optional) - Key algorithm type.

### `multicloud_transit_gateway`
Unified global transit gateway interconnect hub.
- **AWS Target:** `aws_ec2_transit_gateway`
- **GCP Target:** `google_network_connectivity_hub`
- **Azure Target:** `azurerm_virtual_wan` / `azurerm_virtual_hub`
- **Attributes:**
  - `gateway_name` (String, Required) - Gateway hub name.
  - `asn` (Int64, Optional) - BGP Autonomous System Number.
  - `auto_accept_shared_attachments` (Bool, Optional) - Auto accept cross-account attachments.

### `multicloud_private_endpoint`
Unified private link endpoint attachment.
- **AWS Target:** `aws_vpc_endpoint` (Interface Endpoint)
- **GCP Target:** `google_compute_global_forwarding_rule` (PSC)
- **Azure Target:** `azurerm_private_endpoint`
- **Attributes:**
  - `endpoint_name` (String, Required) - Private endpoint name.
  - `vpc_id` (String, Required) - Parent VPC/VNet ID.
  - `service_name` (String, Required) - Target service identifier.
  - `subnet_ids` (List[String], Optional) - Subnet IDs.

### `multicloud_security_center`
Unified cloud security posture and threat monitoring center.
- **AWS Target:** `aws_securityhub_account`
- **GCP Target:** `google_scc_source`
- **Azure Target:** `azurerm_security_center_subscription_pricing`
- **Attributes:**
  - `center_name` (String, Required) - Security center instance name.
  - `tier` (String, Optional) - Tier plan (STANDARD/ADVANCED).
  - `enable_auto_pruning` (Bool, Optional) - Auto-prune resolved security findings.

### `multicloud_kms_policy`
Unified KMS key access policy and grant manager.
- **AWS Target:** `aws_kms_key_policy` / `aws_kms_grant`
- **GCP Target:** `google_kms_crypto_key_iam_binding`
- **Azure Target:** `azurerm_key_vault_access_policy`
- **Attributes:**
  - `policy_name` (String, Required) - Policy name.
  - `key_id` (String, Required) - Target KMS key ID.
  - `policy_json` (String, Required) - Policy JSON document.

### `multicloud_data_pipeline`
Unified ETL batch data processing pipeline.
- **AWS Target:** `aws_glue_crawler` / `aws_emr_cluster`
- **GCP Target:** `google_dataflow_job` / `google_dataproc_cluster`
- **Azure Target:** `azurerm_data_factory_pipeline`
- **Attributes:**
  - `pipeline_name` (String, Required) - Data pipeline job name.
  - `engine` (String, Optional) - Processing engine (SPARK, FLINK, GLUE).
  - `max_workers` (Int64, Optional) - Maximum worker node count.

### `multicloud_feature_store`
Unified MLOps machine learning feature store catalog.
- **AWS Target:** `aws_sagemaker_feature_group`
- **GCP Target:** `google_vertex_ai_featurestore`
- **Azure Target:** `azurerm_machine_learning_workspace`
- **Attributes:**
  - `store_name` (String, Required) - Feature store catalog name.
  - `online_store_enabled` (Bool, Optional) - Enable real-time online feature serving.

### `multicloud_edge_function`
Unified edge serverless PoP code execution function.
- **AWS Target:** `aws_cloudfront_function` / `aws_lambda_function` (Lambda@Edge)
- **GCP Target:** `google_cloudfunctions_function` (CDN Edge)
- **Azure Target:** `azurerm_frontdoor_rules_engine`
- **Attributes:**
  - `function_name` (String, Required) - Edge function name.
  - `runtime` (String, Optional) - Execution runtime (js-1.0, nodejs20).
  - `code_content` (String, Required) - Edge function code content.

### `multicloud_global_anycast_ip`
Unified global anycast static IP address and edge accelerator.
- **AWS Target:** `aws_globalaccelerator_accelerator`
- **GCP Target:** `google_compute_global_address`
- **Azure Target:** `azurerm_traffic_manager_profile`
- **Attributes:**
  - `name` (String, Required) - Global anycast IP identifier.
  - `ip_version` (String, Optional) - `'IPV4'` or `'IPV6'`.
  - `ip_address` (String, Read-Only) - Allocated global anycast IP address.
  - `dns_name` (String, Read-Only) - Global routing DNS name.

### `multicloud_custom_machine_type`
Unified custom vCPU and memory compute specification.
- **AWS Target:** `aws_launch_template`
- **GCP Target:** `google_compute_instance` (`custom-vCPU-RAM`)
- **Azure Target:** `azurerm_linux_virtual_machine`
- **Attributes:**
  - `name` (String, Required) - Custom machine name.
  - `vcpus` (Int64, Required) - Custom vCPU core count.
  - `memory_mb` (Int64, Required) - Custom RAM allocation in MB.
  - `machine_type_spec` (String, Read-Only) - Formatted custom machine type string.

### `multicloud_storage_inventory_report`
Unified object storage bucket inventory and audit report configuration.
- **AWS Target:** `aws_s3_bucket_inventory`
- **GCP Target:** `google_storage_inventory_report_config`
- **Azure Target:** `azurerm_storage_blob_inventory_policy`
- **Attributes:**
  - `name` (String, Required) - Report configuration name.
  - `bucket_name` (String, Required) - Target bucket to audit.
  - `destination_bucket` (String, Required) - Destination bucket for inventory reports.
  - `format` (String, Optional) - `'CSV'`, `'PARQUET'`, or `'ORC'`.
  - `schedule_frequency` (String, Optional) - `'DAILY'` or `'WEEKLY'`.

### `multicloud_workload_identity_pool`
Unified keyless OIDC workload identity federation pool.
- **AWS Target:** `aws_iam_openid_connect_provider`
- **GCP Target:** `google_iam_workload_identity_pool`
- **Azure Target:** `azurerm_federated_identity_credential`
- **Attributes:**
  - `pool_name` (String, Required) - Workload identity pool name.
  - `issuer_url` (String, Required) - OIDC token issuer URL.
  - `allowed_audiences` (List[String], Optional) - Allowed client IDs.
  - `description` (String, Optional) - Pool description.
  - `disabled` (Bool, Optional) - Disable pool state.

### `multicloud_storage_transfer_job`
Unified cross-cloud batch object storage transfer job.
- **AWS Target:** `aws_datasync_task`
- **GCP Target:** `google_storage_transfer_job`
- **Azure Target:** `azurerm_storage_sync`
- **Attributes:**
  - `job_name` (String, Required) - Transfer job name.
  - `source_bucket` (String, Required) - Source bucket name.
  - `destination_bucket` (String, Required) - Destination bucket name.
  - `schedule_start_time` (String, Optional) - ISO-8601 start time.
  - `overwrite_objects` (Bool, Optional) - Overwrite existing destination objects.

### `multicloud_block_volume`
Unified persistent block storage volume.
- **AWS Target:** `aws_ebs_volume`
- **GCP Target:** `google_compute_disk`
- **Azure Target:** `azurerm_managed_disk`
- **Attributes:**
  - `volume_name` (String, Required) - Volume name.
  - `size_gb` (Int64, Required) - Volume size in GB.
  - `volume_type` (String, Optional) - `'ssd'`, `'hdd'`, or `'nvme'`.
  - `iops` (Int64, Optional) - Provisioned IOPS.
  - `encryption_enabled` (Bool, Optional) - Enable disk encryption at rest.

### `multicloud_shared_filesystem`
Unified shared NFS/SMB network filesystem.
- **AWS Target:** `aws_efs_file_system`
- **GCP Target:** `google_filestore_instance`
- **Azure Target:** `azurerm_storage_share`
- **Attributes:**
  - `filesystem_name` (String, Required) - Shared filesystem name.
  - `protocol` (String, Optional) - `'NFSv4'` or `'SMB'`.
  - `performance_mode` (String, Optional) - `'generalPurpose'` or `'maxIO'`.
  - `encryption_enabled` (Bool, Optional) - Enable encryption at rest.

### `multicloud_tls_certificate`
Unified managed X.509 SSL/TLS certificate.
- **AWS Target:** `aws_acm_certificate`
- **GCP Target:** `google_certificate_manager_certificate`
- **Azure Target:** `azurerm_key_vault_certificate`
- **Attributes:**
  - `cert_name` (String, Required) - Certificate name.
  - `domain_name` (String, Required) - Primary FQDN domain name.
  - `validation_method` (String, Optional) - `'DNS'` or `'EMAIL'`.
  - `auto_renew` (Bool, Optional) - Enable automated renewal.

### `multicloud_workflow`
Unified serverless workflow state machine orchestration.
- **AWS Target:** `aws_sfn_state_machine`
- **GCP Target:** `google_workflows_workflow`
- **Azure Target:** `azurerm_logic_app_workflow`
- **Attributes:**
  - `workflow_name` (String, Required) - Workflow name.
  - `definition` (String, Required) - State machine JSON/YAML definition.
  - `workflow_type` (String, Optional) - `'STANDARD'` or `'EXPRESS'`.

### `multicloud_batch_compute`
Unified managed batch compute environment and job pool.
- **AWS Target:** `aws_batch_compute_environment`
- **GCP Target:** `google_batch_job`
- **Azure Target:** `azurerm_batch_pool`
- **Attributes:**
  - `environment_name` (String, Required) - Batch environment name.
  - `max_vcpus` (Int64, Required) - Maximum vCPU capacity.
  - `compute_type` (String, Optional) - `'EC2'`, `'FARGATE'`, or `'SPOT'`.

### `multicloud_backup_vault`
Unified disaster recovery backup vault with optional WORM lock.
- **AWS Target:** `aws_backup_vault`
- **GCP Target:** `google_backup_dr_backup_vault`
- **Azure Target:** `azurerm_data_protection_backup_vault`
- **Attributes:**
  - `vault_name` (String, Required) - Backup vault name.
  - `retention_days` (Int64, Optional) - Retention window in days.
  - `immutable_lock` (Bool, Optional) - Enable WORM lock.
  - `encryption_enabled` (Bool, Optional) - Enable KMS encryption.

### `multicloud_distributed_tracing`
Unified distributed tracing and APM configuration.
- **AWS Target:** `aws_xray_sampling_rule`
- **GCP Target:** `google_cloud_trace_config`
- **Azure Target:** `azurerm_application_insights`
- **Attributes:**
  - `tracing_name` (String, Required) - Tracing rule or workspace name.
  - `sampling_rate` (Float64, Optional) - Sampling ratio (`0.0` to `1.0`).
  - `retention_days` (Int64, Optional) - Trace retention in days.

### `multicloud_budget_alert`
Unified FinOps monthly cost budget and alert threshold policy.
- **AWS Target:** `aws_budgets_budget`
- **GCP Target:** `google_billing_budget`
- **Azure Target:** `azurerm_consumption_budget_subscription`
- **Attributes:**
  - `budget_name` (String, Required) - Budget name.
  - `monthly_limit_usd` (Float64, Required) - Monthly USD spend cap.
  - `alert_threshold_pct` (Int64, Optional) - Alert trigger percentage.
  - `notification_email` (String, Optional) - Alert email recipient.

### `multicloud_vector_index`
Unified AI vector database similarity search index.
- **AWS Target:** `aws_opensearchserverless_collection`
- **GCP Target:** `google_vertex_ai_index`
- **Azure Target:** `azurerm_search_service`
- **Attributes:**
  - `index_name` (String, Required) - Vector index name.
  - `dimensions` (Int64, Required) - Embedding dimensionality (e.g. `1536`).
  - `distance_metric` (String, Optional) - `'COSINE'`, `'DOT_PRODUCT'`, or `'EUCLIDEAN'`.

### `multicloud_service_mesh`
Unified L7 service mesh control plane.
- **AWS Target:** `aws_appmesh_mesh`
- **GCP Target:** `google_network_services_mesh`
- **Azure Target:** `azurerm_kubernetes_fleet_manager`
- **Attributes:**
  - `mesh_name` (String, Required) - Service mesh name.
  - `mtls_mode` (String, Optional) - `'STRICT'` or `'PERMISSIVE'`.
  - `egress_filter` (String, Optional) - `'ALLOW_ALL'` or `'DROP_ALL'`.

---

## 9. Data Sources

### `data "multicloud_resource"`
Reads live or mock state and attributes for an existing cloud resource.
- **Attributes:**
  - `provider_type` (String, Required) - `'aws'`, `'gcp'`, or `'azure'`.
  - `resource_type` (String, Required) - Unified resource type name.
  - `resource_name` (String, Required) - Upstream resource identifier.
  - `region` (String, Optional) - Target cloud region.
  - `status` (String, Read-Only) - Upstream resource status.
  - `attributes` (Map[String], Read-Only) - Upstream resource attributes.

### `data "multicloud_cost_estimate"`
Calculates pre-apply monthly USD costs and Arm64 optimization savings.
- **Attributes:**
  - `provider_type` (String, Required) - `'aws'`, `'gcp'`, or `'azure'`.
  - `resource_type` (String, Required) - Unified resource type name.
  - `size_tier` (String, Optional) - `'small'`, `'medium'`, `'large'`, or explicit SKU.
  - `monthly_cost` (Float64, Read-Only) - Estimated monthly spend in USD.
  - `suggested_tier` (String, Read-Only) - Suggested Arm64 SKU.
  - `estimated_saving` (Float64, Read-Only) - Monthly USD savings.
  - `optimization_note` (String, Read-Only) - FinOps optimization guidance.



