---
# subcategory: "Observability"
page_title: "multicloud_monitoring_dashboard Resource - terraform-provider-multicloud"
description: |-
  Unified Cloud Monitoring Dashboard.
---

# multicloud_monitoring_dashboard (Resource)

Unified Cloud Monitoring Dashboard.

## Cloud Targets
- **AWS Target:** aws_cloudwatch_dashboard
- **GCP Target:** google_monitoring_dashboard
- **Azure Target:** azurerm_portal_dashboard

## How It Works

The `multicloud_monitoring_dashboard` resource provisions visual metrics dashboards across AWS CloudWatch, GCP Cloud Monitoring, and Azure Portal Dashboards.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_monitoring_dashboard" "basic" {
  provider_type  = "aws"
  dashboard_name = "platform-overview"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_monitoring_dashboard" "azure_advanced" {
  provider_type  = "azure"
  dashboard_name = "k8s-metrics-dashboard"

  extra_config = {
    "azure_dashboard_type" = "shared"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `dashboard_name` (String, Required) Dashboard name.
- `dashboard_body` (String, Optional) JSON widget layout.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
