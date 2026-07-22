---
# subcategory: "Security"
page_title: "multicloud_security_group Resource - terraform-provider-multicloud"
description: |-
  Unified Firewall Security Group.
---

# multicloud_security_group (Resource)

Unified Firewall Security Group.

## Cloud Targets
- **AWS Target:** aws_security_group
- **GCP Target:** google_compute_firewall
- **Azure Target:** azurerm_network_security_group

## How It Works

The `multicloud_security_group` resource provisions virtual firewall rules across AWS Security Groups, GCP Compute Firewalls, and Azure Network Security Groups (NSGs).

## Example Usage

### Basic Usage
```hcl
resource "multicloud_security_group" "basic" {
  provider_type = "aws"
  group_name    = "web-sec-group"
  description   = "Allow HTTPS inbound traffic"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_security_group" "gcp_advanced" {
  provider_type = "gcp"
  group_name    = "allow-internal-traffic"
  description   = "GCP Compute internal firewall"

  extra_config = {
    "gcp_target_tags" = "web-node,app-node"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `group_name` (String, Required) Security group name.
- `description` (String, Optional) Group description.
- `vpc_id` (String, Optional) Parent VPC ID.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
