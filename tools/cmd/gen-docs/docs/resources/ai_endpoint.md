---
# subcategory: "Analytics"
page_title: "multicloud_ai_endpoint Resource - terraform-provider-multicloud"
description: |-
  Unified AI / Machine Learning Inference Endpoint.
---

# multicloud_ai_endpoint (Resource)

Unified AI / Machine Learning Inference Endpoint.

## Cloud Targets
- **AWS Target:** aws_sagemaker_endpoint
- **GCP Target:** google_vertex_ai_endpoint
- **Azure Target:** azurerm_cognitive_account

## How It Works

The `multicloud_ai_endpoint` resource provisions real-time inference model deployment endpoints across AWS SageMaker/Bedrock, GCP Vertex AI Endpoints, and Azure OpenAI Service.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_ai_endpoint" "basic" {
  provider_type = "gcp"
  endpoint_name = "llama3-70b-endpoint"
  model_name    = "meta/llama3-70b"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_ai_endpoint" "aws_advanced" {
  provider_type = "aws"
  endpoint_name = "claude-endpoint-aws"
  model_name    = "anthropic.claude-v2"

  extra_config = {
    "aws_instance_count" = "2"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `endpoint_name` (String, Required) Inference endpoint name.
- `model_name` (String, Optional) Deployed model identifier.
- `instance_type` (String, Optional) Compute hardware tier.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
