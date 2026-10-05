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
  value = data.multicloud_cost_estimate.vm_cost.monthly_cost
}

output "arm64_recommendation" {
  value = data.multicloud_cost_estimate.vm_cost.optimization_note
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider (`'aws'`, `'gcp'`, or `'azure'`).
- `resource_type` (String) Unified resource type (e.g., `'virtual_machine'`, `'kubernetes_cluster'`, `'db_instance'`).

### Optional
- `size_tier` (String) Resource size tier (`'small'`, `'medium'`, `'large'`) or explicit SKU (`'m6i.large'`).

### Read-Only
- `id` (String) Cost estimate identifier.
- `monthly_cost` (Float64) Estimated monthly spend in USD.
- `suggested_tier` (String) Recommended Arm64 instance SKU when applicable.
- `estimated_saving` (Float64) Estimated monthly USD savings from Arm64 migration.
- `optimization_note` (String) Actionable FinOps recommendation summary.
