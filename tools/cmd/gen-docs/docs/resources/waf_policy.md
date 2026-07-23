---
# subcategory: "Security"
page_title: "multicloud_waf_policy Resource - terraform-provider-multicloud"
description: |-
  Unified Web Application Firewall Policy.
---

# multicloud_waf_policy (Resource)

Unified Web Application Firewall Policy.

## Cloud Targets
- **AWS Target:** aws_wafv2_web_acl
- **GCP Target:** google_compute_security_policy
- **Azure Target:** azurerm_web_application_firewall_policy

## How It Works

The `multicloud_waf_policy` resource provisions L7 Web Application Firewall rules and bot management across AWS WAFv2, GCP Cloud Armor, and Azure WAF.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_waf_policy" "basic" {
  provider_type  = "aws"
  policy_name    = "app-waf-policy"
  default_action = "allow"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_waf_policy" "gcp_advanced" {
  provider_type  = "gcp"
  policy_name    = "cloud-armor-policy"
  default_action = "allow"

  extra_config = {
    "gcp_adaptive_protection" = "enabled"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `policy_name` (String, Required) WAF policy name.
- `default_action` (String, Optional) 'allow' or 'block'.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
