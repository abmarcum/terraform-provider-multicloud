---
# subcategory: "Networking"
page_title: "multicloud_private_endpoint Resource - terraform-provider-multicloud"
description: |-
  Unified Private Link Endpoint Attachment.
---

# multicloud_private_endpoint (Resource)

Unified Private Link Endpoint Attachment.

## Cloud Targets
- **AWS Target:** aws_vpc_endpoint
- **GCP Target:** google_compute_global_forwarding_rule
- **Azure Target:** azurerm_private_endpoint

## How It Works

The `multicloud_private_endpoint` resource attaches private endpoint interfaces to VPCs/VNets for secure private service access across AWS, GCP, and Azure.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_private_endpoint" "basic" {
  provider_type = "aws"
  endpoint_name = "s3-private-endpoint"
  vpc_id        = multicloud_virtual_network.main.id
  service_name  = "com.amazonaws.us-east-1.s3"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_private_endpoint" "azure_advanced" {
  provider_type = "azure"
  endpoint_name = "kv-private-endpoint"
  vpc_id        = "vnet-id"
  service_name  = "Microsoft.KeyVault"

  extra_config = {
    "azure_private_dns_zone_group" = "default"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `endpoint_name` (String, Required) Private endpoint name.
- `vpc_id` (String, Required) Parent VPC/VNet ID.
- `service_name` (String, Required) Target service identifier.
- `subnet_ids` (List[String], Optional) Placement subnet IDs.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
