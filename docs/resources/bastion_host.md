---
# subcategory: "Security"
page_title: "multicloud_bastion_host Resource - terraform-provider-multicloud"
description: |-
  Unified Managed Bastion Jump Host.
---

# multicloud_bastion_host (Resource)

Unified Managed Bastion Jump Host.

## Cloud Targets
- **AWS Target:** aws_ec2_instance_connect_endpoint
- **GCP Target:** google_iap_tunnel
- **Azure Target:** azurerm_bastion_host

## How It Works

The `multicloud_bastion_host` resource provisions managed SSH jump servers and identity-aware proxy endpoints across AWS EC2 Instance Connect, GCP IAP, and Azure Bastion.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_bastion_host" "basic" {
  provider_type = "aws"
  host_name     = "prod-bastion"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_bastion_host" "azure_advanced" {
  provider_type = "azure"
  host_name     = "azure-bastion"

  extra_config = {
    "azure_sku" = "Standard"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `host_name` (String, Required) Bastion host identifier.
- `vpc_id` (String, Optional) Parent VPC ID.
- `subnet_id` (String, Optional) Target Subnet ID.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
