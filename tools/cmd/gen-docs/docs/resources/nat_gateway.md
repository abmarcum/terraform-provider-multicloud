---
# subcategory: "Network"
page_title: "multicloud_nat_gateway Resource - terraform-provider-multicloud"
description: |-
  Unified NAT Gateway.
---

# multicloud_nat_gateway (Resource)

Unified NAT Gateway.

## Cloud Targets
- **AWS Target:** aws_nat_gateway
- **GCP Target:** google_compute_router_nat
- **Azure Target:** azurerm_nat_gateway

## How It Works

The `multicloud_nat_gateway` resource provisions outbound Network Address Translation (NAT) gateways across AWS NAT Gateways, GCP Cloud NAT, and Azure NAT Gateways to grant private subnet workloads internet egress.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_nat_gateway" "basic" {
  provider_type = "aws"
  gateway_name  = "main-nat-gw"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_nat_gateway" "advanced" {
  provider_type = "aws"
  gateway_name  = "ha-nat-gw"
  subnet_id     = "subnet-abc12345"

  extra_config = {
    "aws_secondary_allocation_ids" = "eipalloc-98765432"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `gateway_name` (String, Required) Gateway name.
- `subnet_id` (String, Optional) Target Subnet ID.
- `allocation_id` (String, Optional) Static IP ID.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
