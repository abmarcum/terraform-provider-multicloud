---
# subcategory: "Storage"
page_title: "multicloud_storage_bucket Resource - terraform-provider-multicloud"
description: |-
  Unified Object Storage Bucket supporting AWS S3, GCP Cloud Storage, and Azure Storage Container.
---

# multicloud_storage_bucket (Resource)

Unified Object Storage Bucket supporting AWS S3, GCP Cloud Storage, and Azure Storage Container.

## Cloud Targets
- **AWS Target:** aws_s3_bucket
- **GCP Target:** google_storage_bucket
- **Azure Target:** azurerm_storage_container

## How It Works

The `multicloud_storage_bucket` resource provisions object storage containers across AWS, GCP, and Azure. It sanitizes bucket names to meet cloud-specific DNS constraints, enables versioning and server-side encryption pre-apply, and maps extra S3/GCS/ARM settings via `extra_config`.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_storage_bucket" "aws_basic" {
  provider_type      = "aws"
  bucket_name        = "my-app-storage-bucket-aws"
  region             = "us-west-2"
  versioning_enabled = true
  encryption_enabled = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_storage_bucket" "gcp_advanced" {
  provider_type      = "gcp"
  bucket_name        = "my-app-storage-bucket-gcp"
  region             = "us-central1"
  versioning_enabled = true
  encryption_enabled = true

  extra_config = {
    "gcp_storage_class" = "NEARLINE"
    "gcp_location_type" = "dual-region"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `bucket_name` (String, Required) Name of the storage bucket.
- `versioning_enabled` (Bool, Optional) Enable object versioning.
- `encryption_enabled` (Bool, Optional) Enable server-side encryption.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
