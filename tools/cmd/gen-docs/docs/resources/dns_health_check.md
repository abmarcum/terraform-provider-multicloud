---
# subcategory: "DNS"
page_title: "multicloud_dns_health_check Resource - terraform-provider-multicloud"
description: |-
  Unified Automated DNS Health Check Monitor.
---

# multicloud_dns_health_check (Resource)

Unified Automated DNS Health Check Monitor.

## Cloud Targets
- **AWS Target:** aws_route53_health_check
- **GCP Target:** google_monitoring_uptime_check_config
- **Azure Target:** azurerm_traffic_manager_endpoint

## How It Works

The `multicloud_dns_health_check` resource provisions endpoint monitoring probes across AWS Route53 Health Checks, GCP Uptime Checks, and Azure Traffic Manager Probes.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dns_health_check" "basic" {
  provider_type     = "aws"
  check_name        = "api-health"
  type              = "HTTPS"
  fqdn              = "api.example.com"
  port              = 443
  resource_path     = "/healthz"
  failure_threshold = 3
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dns_health_check" "azure_advanced" {
  provider_type     = "azure"
  check_name        = "web-probe"
  type              = "HTTP"
  fqdn              = "web.example.com"
  port              = 80
  failure_threshold = 2

  extra_config = {
    "azure_probe_interval" = "30"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `check_name` (String, Required) Health check name.
- `type` (String, Required) Probe protocol (HTTP, HTTPS, TCP).
- `fqdn` (String, Optional) Target FQDN.
- `ip_address` (String, Optional) Target IP.
- `port` (Int64, Optional) Port number.
- `resource_path` (String, Optional) HTTP probe URI path.
- `failure_threshold` (Int64, Optional) Failure threshold count.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
