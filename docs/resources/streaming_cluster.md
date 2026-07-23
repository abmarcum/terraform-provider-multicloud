---
# subcategory: "Messaging"
page_title: "multicloud_streaming_cluster Resource - terraform-provider-multicloud"
description: |-
  Unified Managed Apache Kafka Streaming Cluster.
---

# multicloud_streaming_cluster (Resource)

Unified Managed Apache Kafka Streaming Cluster.

## Cloud Targets
- **AWS Target:** aws_msk_cluster
- **GCP Target:** google_managed_kafka_cluster
- **Azure Target:** azurerm_eventhub_namespace

## How It Works

The `multicloud_streaming_cluster` resource provisions managed event streaming clusters across AWS MSK, GCP Managed Service for Apache Kafka, and Azure Event Hubs.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_streaming_cluster" "basic" {
  provider_type = "aws"
  cluster_name  = "events-kafka"
  kafka_version = "3.5.1"
  node_count    = 3
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_streaming_cluster" "gcp_advanced" {
  provider_type = "gcp"
  cluster_name  = "kafka-gcp-prod"
  kafka_version = "3.4.0"
  node_count    = 6

  extra_config = {
    "gcp_cpu_per_node" = "4"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `cluster_name` (String, Required) Streaming cluster name.
- `kafka_version` (String, Optional) Kafka engine version.
- `node_count` (Int64, Optional) Broker node count.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
