---
# subcategory: "Security"
page_title: "multicloud_secret Resource - terraform-provider-multicloud"
description: |-
  Unified Key-Value Secret Vault Store.
---

# multicloud_secret (Resource)

Unified Key-Value Secret Vault Store.

## Cloud Targets
- **AWS Target:** aws_secretsmanager_secret
- **GCP Target:** google_secret_manager_secret
- **Azure Target:** azurerm_key_vault_secret

## How It Works

The `multicloud_secret` resource provisions encrypted key-value secrets in AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_secret" "basic" {
  provider_type = "gcp"
  secret_name   = "db-password"
  secret_string = var.db_password
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_secret" "aws_advanced" {
  provider_type = "aws"
  secret_name   = "third-party-api-key"
  secret_string = var.api_key

  extra_config = {
    "aws_recovery_window_in_days" = "7"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `secret_name` (String, Required) Secret name.
- `description` (String, Optional) Secret description.
- `secret_string` (String, Optional, Sensitive) Secret payload.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
