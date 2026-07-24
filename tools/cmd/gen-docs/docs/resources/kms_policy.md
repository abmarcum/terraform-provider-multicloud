---
# subcategory: "Security"
page_title: "multicloud_kms_policy Resource - terraform-provider-multicloud"
description: |-
  Unified KMS Key Access Policy & Grant Management.
---

# multicloud_kms_policy (Resource)

Unified KMS Key Access Policy & Grant Management.

## Cloud Targets
- **AWS Target:** aws_kms_key_policy
- **GCP Target:** google_kms_crypto_key_iam_binding
- **Azure Target:** azurerm_key_vault_access_policy

## How It Works

The `multicloud_kms_policy` resource provisions cryptographic key access policies, grant permissions, and IAM key bindings across AWS KMS, GCP KMS, and Azure Key Vault.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_kms_policy" "basic" {
  provider_type = "aws"
  policy_name   = "key-access-policy"
  key_id        = multicloud_kms_key.app_key.id
  policy_json   = jsonencode({ Version = "2012-10-17", Statement = [] })
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_kms_policy" "azure_advanced" {
  provider_type = "azure"
  policy_name   = "kv-policy"
  key_id        = "kv-key-id"
  policy_json   = jsonencode({ key_permissions = ["Get", "List", "Encrypt", "Decrypt"] })

  extra_config = {
    "azure_tenant_id" = "00000000-0000-0000-0000-000000000000"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `policy_name` (String, Required) Policy identifier.
- `key_id` (String, Required) Target KMS Key ID.
- `policy_json` (String, Required) IAM/KMS policy document JSON.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
