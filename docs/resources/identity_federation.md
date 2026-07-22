---
# subcategory: "Security"
page_title: "multicloud_identity_federation Resource - terraform-provider-multicloud"
description: |-
  Multi-cloud OIDC workload identity federation.
---

# multicloud_identity_federation (Resource)

Multi-cloud OIDC workload identity federation.

## Cloud Targets
- **AWS Target:** IAM OIDC Provider
- **GCP Target:** Workload Identity
- **Azure Target:** Entra ID Workload Identity

## How It Works

The `multicloud_identity_federation` resource configures passwordless OIDC workload identity trusts across AWS IAM OIDC, GCP Workload Identity, and Azure Entra ID, allowing external workloads (e.g. Kubernetes, GitHub Actions) to authenticate securely without long-lived static secrets.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_identity_federation" "basic" {
  provider_type   = "aws"
  federation_name = "k8s-oidc"
  issuer_url      = "https://container.googleapis.com/v1/projects/my-proj/locations/us-central1/clusters/my-cluster"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_identity_federation" "azure_advanced" {
  provider_type   = "azure"
  federation_name = "github-actions-federation"
  issuer_url      = "https://token.actions.githubusercontent.com"

  extra_config = {
    "azure_audience" = "api://AzureADTokenExchange"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `federation_name` (String, Required) Federation identifier.
- `issuer_url` (String, Optional) OIDC issuer URL.
- `client_id_list` (List[String], Optional) Target client IDs.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
