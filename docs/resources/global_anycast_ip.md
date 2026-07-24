---
# subcategory: "Networking"
page_title: "multicloud_global_anycast_ip Resource - terraform-provider-multicloud"
description: |-
  Unified Global Anycast IP Address & Accelerator Endpoint.
---

# multicloud_global_anycast_ip (Resource)

Unified Global Anycast IP Address & Accelerator Endpoint.

## Cloud Targets
- **AWS Target:** aws_globalaccelerator_accelerator
- **GCP Target:** google_compute_global_address
- **Azure Target:** azurerm_traffic_manager_profile

## How It Works

The `multicloud_global_anycast_ip` resource provisions global anycast IP routing and low-latency network acceleration across AWS Global Accelerator, GCP Global External Static IP, and Azure Traffic Manager.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_global_anycast_ip" "basic" {
  provider_type   = "aws"
  name            = "global-accelerator-ip"
  ip_address_type = "IPV4"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_global_anycast_ip" "gcp_advanced" {
  provider_type   = "gcp"
  name            = "global-static-ip"
  ip_address_type = "IPV4"

  extra_config = {
    "gcp_network_tier" = "PREMIUM"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `name` (String, Required) Anycast IP resource name.
- `ip_address_type` (String, Optional) Protocol version ('IPV4' or 'IPV6').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
