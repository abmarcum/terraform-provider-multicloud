---
# subcategory: "Compute"
page_title: "multicloud_virtual_machine Resource - terraform-provider-multicloud"
description: |-
  Unified Virtual Machine Compute Instance supporting AWS EC2, GCP Compute Engine, and Azure VMs.
---

# multicloud_virtual_machine (Resource)

Unified Virtual Machine Compute Instance supporting AWS EC2, GCP Compute Engine, and Azure VMs.

## Cloud Targets
- **AWS Target:** aws_instance
- **GCP Target:** google_compute_instance
- **Azure Target:** azurerm_linux_virtual_machine

## How It Works

The `multicloud_virtual_machine` resource maps standardized compute tier profiles ('small', 'medium', 'large') to native instance types (`t3.medium`, `e2-standard-2`, `Standard_B2s`). It manages SSH key deployment, VPC subnet association, and tags across AWS, GCP, and Azure.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_virtual_machine" "gcp_basic" {
  provider_type = "gcp"
  vm_name       = "web-server-gcp"
  size_tier     = "medium"
  region        = "us-central1"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_virtual_machine" "aws_advanced" {
  provider_type = "aws"
  vm_name       = "database-host-aws"
  size_tier     = "large"
  region        = "us-west-2"
  subnet_id     = "subnet-0123456789abcdef0"

  extra_config = {
    "aws_ebs_optimized" = "true"
    "aws_monitoring"    = "detailed"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `vm_name` (String, Required) Virtual machine instance name.
- `size_tier` (String, Optional) Instance tier ('small', 'medium', 'large').
- `image_id` (String, Optional) Custom image ID.
- `subnet_id` (String, Optional) VPC subnet ID.
- `ssh_public_key` (String, Optional, Sensitive) Public SSH key.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
