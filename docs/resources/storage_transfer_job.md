---
# subcategory: "Storage"
page_title: "multicloud_storage_transfer_job Resource - terraform-provider-multicloud"
description: |-
  Unified Cross-Cloud Storage Transfer Job.
---

# multicloud_storage_transfer_job (Resource)

Unified Cross-Cloud Storage Transfer Job.

## Cloud Targets
- **AWS Target:** aws_datasync_task
- **GCP Target:** google_storage_transfer_job
- **Azure Target:** azurerm_storage_sync

## How It Works

The `multicloud_storage_transfer_job` resource provisions batch object migration and data sync jobs across GCP Storage Transfer Service Jobs, AWS DataSync Tasks, and Azure Storage Sync.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_storage_transfer_job" "basic" {
  provider_type      = "gcp"
  job_name           = "aws-to-gcp-sync"
  source_bucket      = "my-aws-s3-bucket"
  destination_bucket = "my-gcp-gcs-bucket"
  overwrite_objects  = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_storage_transfer_job" "aws_advanced" {
  provider_type      = "aws"
  job_name           = "datasync-s3-backup"
  source_bucket      = "prod-data"
  destination_bucket = "dr-backup"
  schedule_start_time= "2026-08-01T00:00:00Z"

  extra_config = {
    "aws_verify_mode" = "POINT_IN_TIME_CONSISTENT"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `job_name` (String, Required) Data transfer job name.
- `source_bucket` (String, Required) Source bucket name or URL.
- `destination_bucket` (String, Required) Destination bucket name or URL.
- `schedule_start_time` (String, Optional) ISO-8601 schedule start time.
- `overwrite_objects` (Bool, Optional) Overwrite destination objects.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
