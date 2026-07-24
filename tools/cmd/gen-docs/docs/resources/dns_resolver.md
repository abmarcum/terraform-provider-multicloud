---
# subcategory: "DNS"
page_title: "multicloud_dns_resolver Resource - terraform-provider-multicloud"
description: |-
  Unified Hybrid Cloud DNS Resolver Endpoint.
---

# multicloud_dns_resolver (Resource)

Unified Hybrid Cloud DNS Resolver Endpoint.

## Cloud Targets
- **AWS Target:** aws_route53_resolver_endpoint
- **GCP Target:** google_dns_policy
- **Azure Target:** azurerm_private_dns_resolver

## How It Works

The `multicloud_dns_resolver` resource provisions inbound/outbound hybrid DNS forwarding endpoints across AWS Route53 Resolver, GCP DNS Policies, and Azure Private DNS Resolver.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_dns_resolver" "basic" {
  provider_type = "aws"
  resolver_name = "inbound-resolver"
  direction     = "INBOUND"
  vpc_id        = multicloud_virtual_network.main.id
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_dns_resolver" "gcp_advanced" {
  provider_type = "gcp"
  resolver_name = "dns-forwarder"
  direction     = "OUTBOUND"
  vpc_id        = "gcp-vpc-id"

  extra_config = {
    "gcp_logging" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `resolver_name` (String, Required) Resolver name.
- `direction` (String, Required) INBOUND or OUTBOUND.
- `vpc_id` (String, Required) Parent VPC ID.
- `ip_configurations` (List[String], Optional) Subnet/IP list.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
