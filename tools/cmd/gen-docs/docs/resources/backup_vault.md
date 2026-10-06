---
# subcategory: "Disaster Recovery"
page_title: "multicloud_backup_vault Resource - terraform-provider-multicloud"
description: |-
  Unified Disaster Recovery Backup Vault supporting AWS Backup, GCP Backup & DR, and Azure Data Protection Backup Vault.
---

# multicloud_backup_vault (Resource)

Unified Disaster Recovery Backup Vault supporting AWS Backup, GCP Backup & DR, and Azure Data Protection Backup Vault.

## Cloud Targets
- **AWS Target:** aws_backup_vault
- **GCP Target:** google_backup_dr_backup_vault
- **Azure Target:** azurerm_data_protection_backup_vault

## How It Works

The `multicloud_backup_vault` resource provisions encrypted, optionally WORM-locked backup vaults across AWS Backup, GCP Backup & DR Service, and Azure Data Protection.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_backup_vault" "basic" {
  provider_type      = "aws"
  vault_name         = "prod-compliance-vault"
  retention_days     = 90
  encryption_enabled = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_backup_vault" "gcp_advanced" {
  provider_type      = "gcp"
  vault_name         = "immutable-dr-vault"
  retention_days     = 365
  immutable_lock     = true
  encryption_enabled = true

  extra_config = {
    "gcp_force_update" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `vault_name` (String, Required) Name of the backup vault.
- `retention_days` (Int64, Optional) Recovery point retention in days.
- `immutable_lock` (Bool, Optional) Enable WORM immutable lock.
- `encryption_enabled` (Bool, Optional) Enable KMS encryption.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
