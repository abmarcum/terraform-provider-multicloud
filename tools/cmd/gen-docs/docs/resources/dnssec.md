---
# subcategory: "DNS"
page_title: "multicloud_dnssec Resource - terraform-provider-multicloud"
description: |-
  Unified DNSSEC Zone Signing Security.
---

# multicloud_dnssec (Resource)

Unified DNSSEC Zone Signing Security.

## Cloud Targets
- **AWS Target:** aws_route53_key_signing_key
- **GCP Target:** google_dns_managed_zone
- **Azure Target:** azurerm_dns_zone

## How It Works

The `multicloud_dnssec` resource provisions DNSSEC key signing and cryptographic zone authentication across AWS Route53 DNSSEC, GCP Cloud DNSSEC, and Azure DNSSEC.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dnssec" "basic" {
  provider_type = "aws"
  zone_id       = multicloud_dns_zone.main.id
  state         = "ON"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dnssec" "gcp_advanced" {
  provider_type = "gcp"
  zone_id       = "pub-zone"
  state         = "ON"
  key_type      = "RSASHA256"

  extra_config = {
    "gcp_ksk_algorithm" = "rsasha256"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `zone_id` (String, Required) Target DNS zone ID.
- `state` (String, Required) Signing state (ON/OFF).
- `key_type` (String, Optional) Key algorithm type.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
