---
# subcategory: "Security"
page_title: "multicloud_secret_rotator Resource - terraform-provider-multicloud"
description: |-
  Multi-cloud automated secret key rotator schedule.
---

# multicloud_secret_rotator (Resource)

Multi-cloud automated secret key rotator schedule.

## Cloud Targets
- **AWS Target:** Secrets Manager Rotation
- **GCP Target:** Secret Manager Rotation
- **Azure Target:** Key Vault Secret Rotation

## How It Works

The `multicloud_secret_rotator` resource configures automated secret rotation schedules in AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault, triggering key generation functions at specified interval days.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_secret_rotator" "basic" {
  provider_type = "aws"
  rotator_name  = "db-key-rotator"
  rotation_days = 90
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_secret_rotator" "advanced" {
  provider_type = "azure"
  rotator_name  = "api-key-rotator-azure"
  rotation_days = 30

  extra_config = {
    "azure_notify_before_expiry" = "P7D"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `rotator_name` (String, Required) Schedule identifier.
- `secret_id` (String, Optional) Target secret ID.
- `rotation_days` (Int64, Optional) Interval in days.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
