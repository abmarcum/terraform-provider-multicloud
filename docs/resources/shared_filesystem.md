---
# subcategory: "Storage"
page_title: "multicloud_shared_filesystem Resource - terraform-provider-multicloud"
description: |-
  Unified Shared Network Filesystem supporting AWS EFS, GCP Cloud Filestore, and Azure Files.
---

# multicloud_shared_filesystem (Resource)

Unified Shared Network Filesystem supporting AWS EFS, GCP Cloud Filestore, and Azure Files.

## Cloud Targets
- **AWS Target:** aws_efs_file_system
- **GCP Target:** google_filestore_instance
- **Azure Target:** azurerm_storage_share

## How It Works

The `multicloud_shared_filesystem` resource provisions elastic NFSv4 and SMB network file shares across AWS Elastic File System (EFS), GCP Cloud Filestore, and Azure Files.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_shared_filesystem" "basic" {
  provider_type      = "aws"
  filesystem_name    = "shared-assets-efs"
  protocol           = "NFSv4"
  encryption_enabled = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_shared_filesystem" "gcp_advanced" {
  provider_type      = "gcp"
  filesystem_name    = "filestore-ml-training"
  performance_mode   = "maxIO"
  encryption_enabled = true

  extra_config = {
    "gcp_tier" = "ENTERPRISE"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `filesystem_name` (String, Required) Name of the shared network filesystem.
- `protocol` (String, Optional) Network file sharing protocol ('NFSv4' or 'SMB').
- `performance_mode` (String, Optional) Throughput mode ('generalPurpose' or 'maxIO').
- `encryption_enabled` (Bool, Optional) Enable encryption at rest.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
