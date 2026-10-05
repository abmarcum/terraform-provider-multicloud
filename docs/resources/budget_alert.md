---
# subcategory: "FinOps"
page_title: "multicloud_budget_alert Resource - terraform-provider-multicloud"
description: |-
  Unified FinOps Cost Budget & Alert Policy supporting AWS Budgets, GCP Billing Budgets, and Azure Consumption Budgets.
---

# multicloud_budget_alert (Resource)

Unified FinOps Cost Budget & Alert Policy supporting AWS Budgets, GCP Billing Budgets, and Azure Consumption Budgets.

## Cloud Targets
- **AWS Target:** aws_budgets_budget
- **GCP Target:** google_billing_budget
- **Azure Target:** azurerm_consumption_budget_subscription

## How It Works

The `multicloud_budget_alert` resource configures monthly spend limits and automated threshold notifications across AWS Budgets, GCP Cloud Billing Budgets, and Azure Consumption Budgets.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_budget_alert" "basic" {
  provider_type       = "aws"
  budget_name         = "platform-monthly-budget"
  monthly_limit_usd   = 5000.00
  alert_threshold_pct = 80
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_budget_alert" "gcp_advanced" {
  provider_type       = "gcp"
  budget_name         = "ai-cluster-spend-cap"
  monthly_limit_usd   = 12000.00
  alert_threshold_pct = 90
  notification_email  = "finops@example.com"

  extra_config = {
    "gcp_credit_types_treatment" = "INCLUDE_ALL_CREDITS"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `budget_name` (String, Required) Name of the monthly cost budget.
- `monthly_limit_usd` (Float64, Required) Monthly spend limit in USD.
- `alert_threshold_pct` (Int64, Optional) Alert threshold percentage.
- `notification_email` (String, Optional) Email recipient for budget alerts.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
