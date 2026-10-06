---
# subcategory: "Data Sources"
page_title: "multicloud_cost_estimate Data Source - terraform-provider-multicloud"
description: |-
  Calculates monthly USD cost estimates and Arm64 Graviton/Ampere/Cobalt optimization savings for any unified resource.
---

# multicloud_cost_estimate (Data Source)

Queries the pre-apply FinOps pricing engine and Arm64 architecture cost optimizer for any unified resource across **AWS**, **GCP**, and **Azure**.

## Example Usage

```hcl
data "multicloud_cost_estimate" "vm_cost" {
  provider_type = "aws"
  resource_type = "virtual_machine"
  size_tier     = "large"
}

output "estimated_monthly_usd" {
  value = data.multicloud_cost_estimate.vm_cost.monthly_cost_usd
}

output "arm64_recommendation" {
  value = data.multicloud_cost_estimate.vm_cost.suggested_tier
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider (`'aws'`, `'gcp'`, or `'azure'`).
- `resource_type` (String) Unified resource type (e.g., `'virtual_machine'`, `'kubernetes_cluster'`, `'db_instance'`).

### Optional
- `size_tier` (String) Instance size tier (`'small'`, `'medium'`, `'large'`).

### Read-Only
- `id` (String) Unique identifier for the cost estimate query.
- `monthly_cost_usd` (Float64) Estimated monthly cost in USD.
- `suggested_tier` (String) Recommended ARM/Graviton/Tau/Ampere instance tier if applicable.
- `estimated_saving_usd` (Float64) Estimated monthly savings in USD when switching to the suggested tier.
