---
# subcategory: "Observability"
page_title: "multicloud_distributed_tracing Resource - terraform-provider-multicloud"
description: |-
  Unified Distributed Tracing & APM supporting AWS X-Ray, GCP Cloud Trace, and Azure Application Insights.
---

# multicloud_distributed_tracing (Resource)

Unified Distributed Tracing & APM supporting AWS X-Ray, GCP Cloud Trace, and Azure Application Insights.

## Cloud Targets
- **AWS Target:** aws_xray_sampling_rule
- **GCP Target:** google_cloud_trace_config
- **Azure Target:** azurerm_application_insights

## How It Works

The `multicloud_distributed_tracing` resource configures end-to-end distributed request tracing and sampling policies across AWS X-Ray, GCP Cloud Trace, and Azure Application Insights.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_distributed_tracing" "basic" {
  provider_type = "aws"
  tracing_name  = "checkout-xray-sampling"
  sampling_rate = 0.10
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_distributed_tracing" "azure_advanced" {
  provider_type  = "azure"
  tracing_name   = "prod-app-insights"
  sampling_rate  = 0.25
  retention_days = 90

  extra_config = {
    "azure_application_type" = "web"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `tracing_name` (String, Required) Name of the tracing rule or APM workspace.
- `sampling_rate` (Float64, Optional) Trace sampling ratio (0.0 to 1.0).
- `retention_days` (Int64, Optional) Span retention duration in days.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
