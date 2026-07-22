---
# subcategory: "Messaging"
page_title: "multicloud_message_queue Resource - terraform-provider-multicloud"
description: |-
  Unified Message Queue.
---

# multicloud_message_queue (Resource)

Unified Message Queue.

## Cloud Targets
- **AWS Target:** aws_sqs_queue
- **GCP Target:** google_pubsub_subscription
- **Azure Target:** azurerm_servicebus_queue

## How It Works

The `multicloud_message_queue` resource provisions asynchronous message queues across AWS SQS, GCP Pub/Sub Subscriptions, and Azure Service Bus Queues.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_message_queue" "basic" {
  provider_type             = "aws"
  queue_name                = "order-processing-queue"
  delay_seconds             = 0
  message_retention_seconds = 86400
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_message_queue" "aws_fifo_advanced" {
  provider_type             = "aws"
  queue_name                = "financial-transactions.fifo"
  message_retention_seconds = 1209600

  extra_config = {
    "aws_fifo_queue"                  = "true"
    "aws_content_based_deduplication" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `queue_name` (String, Required) Queue name.
- `delay_seconds` (Int64, Optional) Delivery delay.
- `max_message_size` (Int64, Optional) Max size in bytes.
- `message_retention_seconds` (Int64, Optional) Retention seconds.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
