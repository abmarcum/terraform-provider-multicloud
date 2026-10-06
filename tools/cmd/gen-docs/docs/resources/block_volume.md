---
# subcategory: "Storage"
page_title: "multicloud_block_volume Resource - terraform-provider-multicloud"
description: |-
  Unified Block Storage Volume supporting AWS EBS, GCP Persistent Disk, and Azure Managed Disk.
---

# multicloud_block_volume (Resource)

Unified Block Storage Volume supporting AWS EBS, GCP Persistent Disk, and Azure Managed Disk.

## Cloud Targets
- **AWS Target:** aws_ebs_volume
- **GCP Target:** google_compute_disk
- **Azure Target:** azurerm_managed_disk

## How It Works

The `multicloud_block_volume` resource provisions persistent block storage volumes across AWS EBS (`gp3`/`io2`), GCP Persistent Disk (`pd-ssd`), and Azure Managed Disks (`Premium_LRS`).

## Example Usage

### Basic Usage
```hcl
resource "multicloud_block_volume" "basic" {
  provider_type      = "aws"
  volume_name        = "app-data-volume"
  size_gb            = 100
  encryption_enabled = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_block_volume" "gcp_advanced" {
  provider_type      = "gcp"
  volume_name        = "db-nvme-disk"
  size_gb            = 500
  volume_type        = "ssd"
  iops               = 15000
  encryption_enabled = true

  extra_config = {
    "gcp_physical_block_size_bytes" = "4096"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `volume_name` (String, Required) Name of the block storage volume.
- `size_gb` (Int64, Required) Volume capacity in gigabytes (GB).
- `volume_type` (String, Optional) Storage performance tier ('ssd', 'hdd', 'nvme').
- `iops` (Int64, Optional) Provisioned IOPS.
- `encryption_enabled` (Bool, Optional) Enable disk encryption at rest.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
