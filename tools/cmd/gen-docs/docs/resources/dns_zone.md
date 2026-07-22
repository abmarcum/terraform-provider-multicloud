---
# subcategory: "DNS"
page_title: "multicloud_dns_zone Resource - terraform-provider-multicloud"
description: |-
  Unified Managed DNS Zone.
---

# multicloud_dns_zone (Resource)

Unified Managed DNS Zone.

## Cloud Targets
- **AWS Target:** aws_route53_zone
- **GCP Target:** google_dns_managed_zone
- **Azure Target:** azurerm_dns_zone

## How It Works

The `multicloud_dns_zone` resource provisions authoritative DNS zones across AWS Route 53, GCP Cloud DNS, and Azure DNS.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dns_zone" "basic" {
  provider_type = "aws"
  zone_name     = "prod-zone"
  domain_name   = "example.com"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dns_zone" "gcp_advanced" {
  provider_type = "gcp"
  zone_name     = "corp-internal-dns"
  domain_name   = "internal.company.net"

  extra_config = {
    "gcp_visibility" = "private"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `zone_name` (String, Required) DNS zone name.
- `domain_name` (String, Optional) FQDN domain name.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
