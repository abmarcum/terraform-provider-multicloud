---
# subcategory: "Network"
page_title: "multicloud_static_ip Resource - terraform-provider-multicloud"
description: |-
  Unified Static Elastic IP Address.
---

# multicloud_static_ip (Resource)

Unified Static Elastic IP Address.

## Cloud Targets
- **AWS Target:** aws_eip
- **GCP Target:** google_compute_address
- **Azure Target:** azurerm_public_ip

## How It Works

The `multicloud_static_ip` resource allocates static public IPv4 addresses across AWS Elastic IPs, GCP Compute Addresses, and Azure Public IPs.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_static_ip" "basic" {
  provider_type   = "azure"
  ip_name         = "app-public-ip"
  allocation_type = "static"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_static_ip" "gcp_advanced" {
  provider_type   = "gcp"
  ip_name         = "lb-static-ip-gcp"
  allocation_type = "static"

  extra_config = {
    "gcp_network_tier" = "PREMIUM"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `ip_name` (String, Required) IP name.
- `allocation_type` (String, Optional) 'static' or 'dynamic'.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
