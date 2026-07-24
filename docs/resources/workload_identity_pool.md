---
# subcategory: "IAM & Identity"
page_title: "multicloud_workload_identity_pool Resource - terraform-provider-multicloud"
description: |-
  Unified Workload Identity Federation Pool.
---

# multicloud_workload_identity_pool (Resource)

Unified Workload Identity Federation Pool.

## Cloud Targets
- **AWS Target:** aws_iam_openid_connect_provider
- **GCP Target:** google_iam_workload_identity_pool
- **Azure Target:** azurerm_federated_identity_credential

## How It Works

The `multicloud_workload_identity_pool` resource configures keyless OIDC identity federation pools across GCP Workload Identity Pools, AWS IAM OIDC Providers, and Azure Entra ID Federated Credentials.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_workload_identity_pool" "basic" {
  provider_type = "gcp"
  pool_name     = "github-actions-pool"
  issuer_url    = "https://token.actions.githubusercontent.com"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_workload_identity_pool" "aws_advanced" {
  provider_type     = "aws"
  pool_name         = "gitlab-ci-oidc"
  issuer_url        = "https://gitlab.com"
  allowed_audiences = ["https://gitlab.com"]

  extra_config = {
    "aws_thumbprint_list" = "9e99a48a9960b14926bb7f3b02e22da2b0ab7280"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `pool_name` (String, Required) Workload identity pool name.
- `issuer_url` (String, Required) OIDC issuer URL.
- `allowed_audiences` (List, Optional) Allowed audience client IDs.
- `description` (String, Optional) Pool description.
- `disabled` (Bool, Optional) Disable pool state.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
