---
# subcategory: "Network"
page_title: "multicloud_vpc_peering Resource - terraform-provider-multicloud"
description: |-
  Unified Virtual Private Network Peering Connection.
---

# multicloud_vpc_peering (Resource)

Unified Virtual Private Network Peering Connection.

## Cloud Targets
- **AWS Target:** aws_vpc_peering_connection
- **GCP Target:** google_compute_network_peering
- **Azure Target:** azurerm_virtual_network_peering

## How It Works

The `multicloud_vpc_peering` resource provisions non-transitive VPC and VNet peering connections across AWS VPC Peering, GCP Network Peering, and Azure VNet Peering.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_vpc_peering" "basic" {
  provider_type = "aws"
  peering_name  = "vpc-peering-main"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_vpc_peering" "azure_advanced" {
  provider_type = "azure"
  peering_name  = "vnet-peering-hub-spoke"

  extra_config = {
    "azure_allow_forwarded_traffic" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `peering_name` (String, Required) Peering connection name.
- `vpc_id` (String, Optional) Source VPC ID.
- `peer_vpc_id` (String, Optional) Destination VPC ID.
- `peer_region` (String, Optional) Destination VPC region.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
