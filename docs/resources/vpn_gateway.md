---
# subcategory: "Network"
page_title: "multicloud_vpn_gateway Resource - terraform-provider-multicloud"
description: |-
  Unified Virtual Private Network (VPN) Gateway.
---

# multicloud_vpn_gateway (Resource)

Unified Virtual Private Network (VPN) Gateway.

## Cloud Targets
- **AWS Target:** aws_vpn_gateway
- **GCP Target:** google_compute_vpn_gateway
- **Azure Target:** azurerm_virtual_network_gateway

## How It Works

The `multicloud_vpn_gateway` resource provisions encrypted IPsec VPN gateways across AWS VPN Gateways, GCP Cloud VPN, and Azure VNet Gateways to connect hybrid on-premises sites to cloud networks.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_vpn_gateway" "basic" {
  provider_type = "aws"
  gateway_name  = "office-vpn"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_vpn_gateway" "advanced" {
  provider_type = "gcp"
  gateway_name  = "hybrid-cloud-vpn"
  tunnel_ip     = "203.0.113.50"

  extra_config = {
    "gcp_ike_version" = "2"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `gateway_name` (String, Required) Gateway name.
- `vpc_id` (String, Optional) Target VPC ID.
- `tunnel_ip` (String, Optional) Remote tunnel IP.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
