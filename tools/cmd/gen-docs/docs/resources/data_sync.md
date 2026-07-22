---
# subcategory: "Replication"
page_title: "multicloud_data_sync Resource - terraform-provider-multicloud"
description: |-
  Multi-cloud object storage background data sync.
---

# multicloud_data_sync (Resource)

Multi-cloud object storage background data sync.

## Cloud Targets
- **AWS Target:** S3 Cross-Region Replication
- **GCP Target:** Storage Transfer Service
- **Azure Target:** Azure Storage Sync

## How It Works

The `multicloud_data_sync` resource orchestrates cloud-native background object replication jobs between source and destination storage buckets across cloud boundaries using S3 Replication, GCP Storage Transfer Service, or Azure Storage Sync.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_data_sync" "basic" {
  provider_type      = "aws"
  sync_name          = "s3-to-gcs-sync"
  source_bucket      = "s3://my-source-bucket"
  destination_bucket = "gs://my-dest-bucket"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_data_sync" "advanced" {
  provider_type      = "gcp"
  sync_name          = "nightly-analytics-sync"
  source_bucket      = "gs://prod-logs-bucket"
  destination_bucket = "s3://analytics-archive-bucket"
  schedule_cron      = "0 2 * * *"

  extra_config = {
    "gcp_overwrite_objects" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `sync_name` (String, Required) Task identifier.
- `source_bucket` (String, Optional) Source URI.
- `destination_bucket` (String, Optional) Destination URI.
- `schedule_cron` (String, Optional) Cron schedule.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
