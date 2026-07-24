---
# subcategory: "Compute"
page_title: "multicloud_custom_machine_type Resource - terraform-provider-multicloud"
description: |-
  Unified Custom vCPU & Memory Machine Type Allocation.
---

# multicloud_custom_machine_type (Resource)

Unified Custom vCPU & Memory Machine Type Allocation.

## Cloud Targets
- **AWS Target:** aws_launch_template
- **GCP Target:** google_compute_instance (custom-vCPU-RAM)
- **Azure Target:** azurerm_linux_virtual_machine

## How It Works

The `multicloud_custom_machine_type` resource provisions customized vCPU and RAM hardware allocations (e.g. custom-4-16384) across GCP Custom Machine Types, AWS EC2 Launch Templates, and Azure Custom VM sizes.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_custom_machine_type" "basic" {
  provider_type = "gcp"
  name          = "custom-workload"
  vcpus         = 4
  memory_mb     = 16384
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_custom_machine_type" "aws_advanced" {
  provider_type = "aws"
  name          = "custom-ec2-spec"
  vcpus         = 8
  memory_mb     = 32768

  extra_config = {
    "aws_ebs_optimized" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `name` (String, Required) Custom machine name.
- `vcpus` (Int64, Required) Custom vCPU core count.
- `memory_mb` (Int64, Required) Custom RAM memory in megabytes.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
