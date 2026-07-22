---
# subcategory: "Messaging"
page_title: "multicloud_event_bridge Resource - terraform-provider-multicloud"
description: |-
  Unified Event Router & Trigger Bus.
---

# multicloud_event_bridge (Resource)

Unified Event Router & Trigger Bus.

## Cloud Targets
- **AWS Target:** aws_cloudwatch_event_bus
- **GCP Target:** google_eventarc_trigger
- **Azure Target:** azurerm_eventgrid_system_topic

## How It Works

The `multicloud_event_bridge` resource provisions event routing buses across AWS EventBridge, GCP Eventarc, and Azure Event Grid System Topics.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_event_bridge" "basic" {
  provider_type = "aws"
  bus_name      = "custom-app-eventbus"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_event_bridge" "gcp_advanced" {
  provider_type = "gcp"
  bus_name      = "cloud-storage-trigger"

  extra_config = {
    "gcp_destination_run_service" = "event-processor"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `bus_name` (String, Required) Event bus router name.
- `event_source` (String, Optional) Event source identifier.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
