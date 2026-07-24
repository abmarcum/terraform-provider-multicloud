---
# subcategory: "Analytics"
page_title: "multicloud_feature_store Resource - terraform-provider-multicloud"
description: |-
  Unified MLOps Machine Learning Feature Store.
---

# multicloud_feature_store (Resource)

Unified MLOps Machine Learning Feature Store.

## Cloud Targets
- **AWS Target:** aws_sagemaker_feature_group
- **GCP Target:** google_vertex_ai_featurestore
- **Azure Target:** azurerm_machine_learning_workspace

## How It Works

The `multicloud_feature_store` resource provisions centralized MLOps machine learning feature catalogs and low-latency online feature serving across AWS SageMaker Feature Store, GCP Vertex AI Featurestore, and Azure Machine Learning Workspace.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_feature_store" "basic" {
  provider_type        = "aws"
  store_name           = "customer-churn-features"
  online_store_enabled = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_feature_store" "gcp_advanced" {
  provider_type        = "gcp"
  store_name           = "recommendation-features"
  online_store_enabled = true

  extra_config = {
    "gcp_fixed_node_count" = "2"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `store_name` (String, Required) Feature store catalog name.
- `online_store_enabled` (Bool, Optional) Enable real-time low-latency online serving.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
