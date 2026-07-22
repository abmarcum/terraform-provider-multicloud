---
# subcategory: "IAM"
page_title: "multicloud_iam_role Resource - terraform-provider-multicloud"
description: |-
  Unified IAM Role and Identity Policy.
---

# multicloud_iam_role (Resource)

Unified IAM Role and Identity Policy.

## Cloud Targets
- **AWS Target:** aws_iam_role
- **GCP Target:** google_service_account
- **Azure Target:** azurerm_user_assigned_identity

## How It Works

The `multicloud_iam_role` resource provisions identity roles across AWS IAM Roles, GCP Service Accounts, and Azure Managed Identities.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_iam_role" "basic" {
  provider_type = "aws"
  role_name     = "app-execution-role"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_iam_role" "aws_advanced" {
  provider_type       = "aws"
  role_name           = "ecs-task-execution-role"
  assume_role_policy  = file("${path.module}/policies/trust_policy.json")

  extra_config = {
    "aws_max_session_duration" = "7200"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `role_name` (String, Required) Role name.
- `assume_role_policy` (String, Optional) Trust policy JSON.
- `description` (String, Optional) Role description.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
