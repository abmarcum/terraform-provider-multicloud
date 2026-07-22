---
# subcategory: "Messaging"
page_title: "multicloud_pubsub_topic Resource - terraform-provider-multicloud"
description: |-
  Unified Publish/Subscribe Topic.
---

# multicloud_pubsub_topic (Resource)

Unified Publish/Subscribe Topic.

## Cloud Targets
- **AWS Target:** aws_sns_topic
- **GCP Target:** google_pubsub_topic
- **Azure Target:** azurerm_servicebus_topic

## How It Works

The `multicloud_pubsub_topic` resource provisions pub/sub messaging topics across AWS SNS, GCP Pub/Sub, and Azure Service Bus Topics.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_pubsub_topic" "basic" {
  provider_type = "gcp"
  topic_name    = "user-signup-events"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_pubsub_topic" "aws_advanced" {
  provider_type = "aws"
  topic_name    = "payment-alerts-sns"
  display_name  = "Payment Processing Alerts"

  extra_config = {
    "aws_kms_master_key_id" = "alias/aws/sns"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `topic_name` (String, Required) Topic name.
- `display_name` (String, Optional) Friendly label.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
