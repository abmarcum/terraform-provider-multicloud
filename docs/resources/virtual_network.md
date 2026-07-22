---
# subcategory: "Network"
page_title: "multicloud_virtual_network Resource - terraform-provider-multicloud"
description: |-
  Unified Virtual Network (VPC / VNet).
---

# multicloud_virtual_network (Resource)

Unified Virtual Network (VPC / VNet).

## Cloud Targets
- **AWS Target:** aws_vpc
- **GCP Target:** google_compute_network
- **Azure Target:** azurerm_virtual_network

## How It Works

The `multicloud_virtual_network` resource provisions isolated virtual private networks across AWS VPCs, GCP VPC Networks, and Azure VNets, configuring CIDR address spaces and network routing boundaries.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_virtual_network" "basic" {
  provider_type = "aws"
  network_name  = "prod-vpc"
  cidr_block    = "10.0.0.0/16"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_virtual_network" "azure_advanced" {
  provider_type = "azure"
  network_name  = "corp-vnet-azure"
  cidr_block    = "172.16.0.0/12"

  extra_config = {
    "azure_enable_ddos_protection" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `network_name` (String, Required) Network identifier.
- `cidr_block` (String, Optional) IPv4 CIDR range.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
