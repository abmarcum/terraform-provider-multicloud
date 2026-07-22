---
# subcategory: "Network"
page_title: "multicloud_cdn_distribution Resource - terraform-provider-multicloud"
description: |-
  Unified Content Delivery Network (CDN) Distribution.
---

# multicloud_cdn_distribution (Resource)

Unified Content Delivery Network (CDN) Distribution.

## Cloud Targets
- **AWS Target:** aws_cloudfront_distribution
- **GCP Target:** google_compute_backend_service
- **Azure Target:** azurerm_cdn_endpoint

## How It Works

The `multicloud_cdn_distribution` resource provisions edge caching content delivery distributions across AWS CloudFront, GCP Cloud CDN, and Azure CDN Endpoints.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_cdn_distribution" "basic" {
  provider_type     = "aws"
  distribution_name = "assets-cdn"
  origin_domain     = "my-bucket.s3.amazonaws.com"
  enabled           = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_cdn_distribution" "advanced" {
  provider_type     = "aws"
  distribution_name = "global-media-cdn"
  origin_domain     = "media.example.com"
  enabled           = true

  extra_config = {
    "aws_price_class" = "PriceClass_100"
    "aws_ipv6_enabled" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `distribution_name` (String, Required) CDN name.
- `origin_domain` (String, Optional) Origin host domain.
- `enabled` (Bool, Optional) CDN status.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
