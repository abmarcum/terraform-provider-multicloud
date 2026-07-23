---
# subcategory: "Observability"
page_title: "multicloud_log_workspace Resource - terraform-provider-multicloud"
description: |-
  Unified Log Analytics Workspace.
---

# multicloud_log_workspace (Resource)

Unified Log Analytics Workspace.

## Cloud Targets
- **AWS Target:** aws_cloudwatch_log_group
- **GCP Target:** google_logging_project_sink
- **Azure Target:** azurerm_log_analytics_workspace

## How It Works

The `multicloud_log_workspace` resource provisions log aggregation workspaces and sinks across AWS CloudWatch Log Groups, GCP Logging Sinks, and Azure Log Analytics Workspaces.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_log_workspace" "basic" {
  provider_type  = "aws"
  workspace_name = "app-logs"
  retention_days = 30
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_log_workspace" "gcp_advanced" {
  provider_type  = "gcp"
  workspace_name = "gcp-audit-sink"
  retention_days = 365

  extra_config = {
    "gcp_unique_writer_identity" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `workspace_name` (String, Required) Workspace name.
- `retention_days` (Int64, Optional) Retention period in days.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
