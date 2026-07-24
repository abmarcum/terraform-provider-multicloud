---
# subcategory: "DNS"
page_title: "multicloud_dns_record Resource - terraform-provider-multicloud"
description: |-
  Unified DNS Resource Record Set.
---

# multicloud_dns_record (Resource)

Unified DNS Resource Record Set.

## Cloud Targets
- **AWS Target:** aws_route53_record
- **GCP Target:** google_dns_record_set
- **Azure Target:** azurerm_dns_a_record

## How It Works

The `multicloud_dns_record` resource provisions A, AAAA, CNAME, TXT, MX, and NS record sets across AWS Route53, GCP Cloud DNS, and Azure DNS.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dns_record" "basic" {
  provider_type = "aws"
  zone_id       = multicloud_dns_zone.main.id
  record_name   = "api.example.com"
  record_type   = "A"
  ttl           = 300
  records       = ["192.0.2.1"]
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dns_record" "gcp_advanced" {
  provider_type = "gcp"
  zone_id       = "prod-zone"
  record_name   = "app.example.com"
  record_type   = "CNAME"
  ttl           = 60
  records       = ["lb.example.com"]

  extra_config = {
    "gcp_routing_policy" = "WRR"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `zone_id` (String, Required) Parent DNS zone ID.
- `record_name` (String, Required) Record hostname.
- `record_type` (String, Required) DNS type (A, CNAME, etc.).
- `ttl` (Int64, Optional) Time-To-Live.
- `records` (List[String], Optional) List of target IPs/hostnames.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
