---
# subcategory: "Disaster Recovery"
page_title: "multicloud_failover_policy Resource - terraform-provider-multicloud"
description: |-
  Unified Active-Passive Disaster Recovery Failover Policy.
---

# multicloud_failover_policy (Resource)

Unified Active-Passive Disaster Recovery Failover Policy.

## Cloud Targets
- **AWS Target:** Route53 Failover Routing
- **GCP Target:** Cloud DNS Failover Policy
- **Azure Target:** Azure Traffic Manager

## How It Works

The `multicloud_failover_policy` resource configures global active-passive routing between primary and secondary clouds. During runtime, managed cloud DNS edge networks (AWS Route 53, GCP Cloud DNS, Azure Traffic Manager) probe `health_check_url` every 10-30s. If the primary cloud fails, traffic is automatically diverted to the failover cloud without local agent dependencies.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_failover_policy" "basic" {
  policy_name      = "global-app-failover"
  primary_cloud    = "aws"
  failover_cloud   = "gcp"
  health_check_url = "https://app.example.com/healthz"
  auto_failover    = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_failover_policy" "azure_advanced" {
  policy_name      = "enterprise-multi-cloud-dr"
  primary_cloud    = "azure"
  failover_cloud   = "aws"
  health_check_url = "https://api.example.com/healthz"
  auto_failover    = true

  extra_config = {
    "azure_routing_method" = "Performance"
    "azure_ttl"            = "60"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `policy_name` (String, Required) Policy identifier.
- `primary_cloud` (String, Required) Primary cloud provider ('aws', 'gcp', 'azure').
- `failover_cloud` (String, Required) Backup cloud provider.
- `health_check_url` (String, Optional) Endpoint URL.
- `auto_failover` (Bool, Optional) Automated failover status.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
