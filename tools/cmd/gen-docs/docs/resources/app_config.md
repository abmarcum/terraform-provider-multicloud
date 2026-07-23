---
# subcategory: "Security"
page_title: "multicloud_app_config Resource - terraform-provider-multicloud"
description: |-
  Unified Application Configuration Store.
---

# multicloud_app_config (Resource)

Unified Application Configuration Store.

## Cloud Targets
- **AWS Target:** aws_ssm_parameter
- **GCP Target:** google_runtimeconfig_config
- **Azure Target:** azurerm_app_configuration

## How It Works

The `multicloud_app_config` resource provisions runtime key-value application configurations across AWS SSM Parameter Store, GCP Runtime Config, and Azure App Configuration.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_app_config" "basic" {
  provider_type = "aws"
  config_name   = "db-max-connections"
  config_key    = "/app/db/max_connections"
  config_value  = "100"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_app_config" "gcp_advanced" {
  provider_type = "gcp"
  config_name   = "feature-flags"
  config_key    = "enable_new_ui"
  config_value  = "true"

  extra_config = {
    "gcp_label" = "prod"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `config_name` (String, Required) Configuration name.
- `config_key` (String, Optional) Parameter key.
- `config_value` (String, Optional) Parameter value.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
