---
# subcategory: "Security"
page_title: "multicloud_kms_key Resource - terraform-provider-multicloud"
description: |-
  Unified KMS Encryption Key.
---

# multicloud_kms_key (Resource)

Unified KMS Encryption Key.

## Cloud Targets
- **AWS Target:** aws_kms_key
- **GCP Target:** google_kms_crypto_key
- **Azure Target:** azurerm_key_vault_key

## How It Works

The `multicloud_kms_key` resource provisions Customer-Managed KMS Encryption Keys (CMKs) across AWS KMS, GCP KMS, and Azure Key Vault, managing deletion windows and automated key rotation.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_kms_key" "basic" {
  provider_type           = "aws"
  key_name                = "app-kms-key"
  deletion_window_in_days = 30
  enable_key_rotation     = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_kms_key" "gcp_advanced" {
  provider_type       = "gcp"
  key_name            = "gcp-cmk-key"
  enable_key_rotation = true

  extra_config = {
    "gcp_rotation_period" = "7776000s"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `key_name` (String, Required) Key alias name.
- `description` (String, Optional) Key description.
- `deletion_window_in_days` (Int64, Optional) Deletion window.
- `enable_key_rotation` (Bool, Optional) Enable key rotation.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
