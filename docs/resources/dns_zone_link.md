---
# subcategory: "DNS"
page_title: "multicloud_dns_zone_link Resource - terraform-provider-multicloud"
description: |-
  Unified Private DNS Network Association.
---

# multicloud_dns_zone_link (Resource)

Unified Private DNS Network Association.

## Cloud Targets
- **AWS Target:** aws_route53_zone_association
- **GCP Target:** google_dns_managed_zone
- **Azure Target:** azurerm_private_dns_zone_virtual_network_link

## How It Works

The `multicloud_dns_zone_link` resource provisions network attachments linking private DNS zones to VPCs/VNets across AWS, GCP, and Azure.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dns_zone_link" "basic" {
  provider_type = "aws"
  link_name     = "private-link"
  zone_id       = multicloud_dns_zone.private.id
  vpc_id        = multicloud_virtual_network.main.id
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dns_zone_link" "azure_advanced" {
  provider_type        = "azure"
  link_name            = "vnet-link"
  zone_id              = "private-dns-zone-id"
  vpc_id               = "vnet-id"
  registration_enabled = true

  extra_config = {
    "azure_auto_registration" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `link_name` (String, Required) Association name.
- `zone_id` (String, Required) Private DNS zone ID.
- `vpc_id` (String, Required) Target VPC/VNet ID.
- `registration_enabled` (Bool, Optional) Auto-register VM hostnames.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
