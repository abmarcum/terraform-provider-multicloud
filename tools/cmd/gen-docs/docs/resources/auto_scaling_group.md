---
# subcategory: "Compute"
page_title: "multicloud_auto_scaling_group Resource - terraform-provider-multicloud"
description: |-
  Unified Compute Auto Scaling Group.
---

# multicloud_auto_scaling_group (Resource)

Unified Compute Auto Scaling Group.

## Cloud Targets
- **AWS Target:** aws_autoscaling_group
- **GCP Target:** google_compute_instance_group_manager
- **Azure Target:** azurerm_linux_virtual_machine_scale_set

## How It Works

The `multicloud_auto_scaling_group` resource provisions dynamic compute scale groups across AWS ASGs, GCP Instance Group Managers, and Azure VM Scale Sets, managing minimum, maximum, and desired node capacities.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_auto_scaling_group" "basic" {
  provider_type    = "aws"
  group_name       = "web-asg"
  min_size         = 2
  max_size         = 10
  desired_capacity = 4
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_auto_scaling_group" "gcp_advanced" {
  provider_type    = "gcp"
  group_name       = "frontend-mig-gcp"
  min_size         = 3
  max_size         = 15
  desired_capacity = 5

  extra_config = {
    "gcp_target_cpu_utilization" = "0.75"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `group_name` (String, Required) Group name.
- `min_size` (Int64, Optional) Minimum capacity.
- `max_size` (Int64, Optional) Maximum capacity.
- `desired_capacity` (Int64, Optional) Desired capacity.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
