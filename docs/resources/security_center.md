---
# subcategory: "Security"
page_title: "multicloud_security_center Resource - terraform-provider-multicloud"
description: |-
  Unified Cloud Security Posture & Threat Monitoring Baseline.
---

# multicloud_security_center (Resource)

Unified Cloud Security Posture & Threat Monitoring Baseline.

## Cloud Targets
- **AWS Target:** aws_securityhub_account
- **GCP Target:** google_scc_source
- **Azure Target:** azurerm_security_center_subscription_pricing

## How It Works

The `multicloud_security_center` resource enables unified security posture monitoring, compliance auditing, and threat detection across AWS Security Hub, GCP Security Command Center, and Azure Microsoft Defender for Cloud.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_security_center" "basic" {
  provider_type = "aws"
  center_name   = "default-security-center"
  tier          = "STANDARD"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_security_center" "gcp_advanced" {
  provider_type = "gcp"
  center_name   = "org-scc"
  tier          = "ADVANCED"

  extra_config = {
    "gcp_scc_mode" = "PREMIUM"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `center_name` (String, Required) Security posture instance name.
- `tier` (String, Optional) Monitoring plan tier (STANDARD/ADVANCED).
- `enable_auto_pruning` (Bool, Optional) Auto-prune resolved security findings.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
