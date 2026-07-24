---
# subcategory: "Storage"
page_title: "multicloud_storage_inventory_report Resource - terraform-provider-multicloud"
description: |-
  Unified Storage Bucket Inventory & Audit Policy.
---

# multicloud_storage_inventory_report (Resource)

Unified Storage Bucket Inventory & Audit Policy.

## Cloud Targets
- **AWS Target:** aws_s3_bucket_inventory
- **GCP Target:** google_storage_inventory_report_config
- **Azure Target:** azurerm_storage_blob_inventory_policy

## How It Works

The `multicloud_storage_inventory_report` resource configures automated object storage auditing and inventory file generation across AWS S3 Bucket Inventory, GCP Cloud Storage Inventory Reports, and Azure Blob Storage Inventory.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_storage_inventory_report" "basic" {
  provider_type      = "aws"
  name               = "s3-daily-audit"
  bucket_name        = "prod-data-bucket"
  destination_bucket = "prod-inventory-reports"
  format             = "CSV"
  schedule_frequency = "DAILY"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_storage_inventory_report" "gcp_advanced" {
  provider_type      = "gcp"
  name               = "gcs-weekly-inventory"
  bucket_name        = "company-analytics"
  destination_bucket = "company-audit-logs"
  format             = "PARQUET"
  schedule_frequency = "WEEKLY"

  extra_config = {
    "gcp_include_prefixes" = "logs/"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `name` (String, Required) Report configuration name.
- `bucket_name` (String, Required) Target bucket to audit.
- `destination_bucket` (String, Required) Destination report delivery bucket.
- `format` (String, Optional) Report format ('CSV', 'PARQUET', or 'ORC').
- `schedule_frequency` (String, Optional) Frequency ('DAILY' or 'WEEKLY').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
