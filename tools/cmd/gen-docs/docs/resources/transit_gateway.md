---
# subcategory: "Networking"
page_title: "multicloud_transit_gateway Resource - terraform-provider-multicloud"
description: |-
  Unified Global Transit Gateway Interconnect Hub.
---

# multicloud_transit_gateway (Resource)

Unified Global Transit Gateway Interconnect Hub.

## Cloud Targets
- **AWS Target:** aws_ec2_transit_gateway
- **GCP Target:** google_network_connectivity_hub
- **Azure Target:** azurerm_virtual_wan

## How It Works

The `multicloud_transit_gateway` resource provisions multi-region, multi-VPC transit interconnect hubs across AWS Transit Gateway, GCP Network Connectivity Center, and Azure Virtual WAN.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_transit_gateway" "basic" {
  provider_type = "aws"
  gateway_name  = "global-hub"
  asn           = 64512
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_transit_gateway" "azure_advanced" {
  provider_type = "azure"
  gateway_name  = "vwan-hub"

  extra_config = {
    "azure_vwan_type" = "Standard"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `gateway_name` (String, Required) Gateway hub name.
- `asn` (Int64, Optional) BGP Autonomous System Number.
- `auto_accept_shared_attachments` (Bool, Optional) Auto accept cross-account attachments.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
