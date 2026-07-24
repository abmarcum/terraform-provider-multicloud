---
# subcategory: "Compute"
page_title: "multicloud_edge_function Resource - terraform-provider-multicloud"
description: |-
  Unified Edge PoP Serverless Code Execution.
---

# multicloud_edge_function (Resource)

Unified Edge PoP Serverless Code Execution.

## Cloud Targets
- **AWS Target:** aws_cloudfront_function
- **GCP Target:** google_cloudfunctions_function
- **Azure Target:** azurerm_frontdoor_rules_engine

## How It Works

The `multicloud_edge_function` resource executes ultra-low-latency serverless code at CDN Points of Presence (PoPs) across AWS CloudFront Functions / Lambda@Edge, GCP Cloud Functions Edge, and Azure Front Door Rules Engine.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_edge_function" "basic" {
  provider_type = "aws"
  function_name = "url-rewrite"
  runtime       = "cloudfront-js-1.0"
  code_content  = "function handler(event) { return event.request; }"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_edge_function" "azure_advanced" {
  provider_type = "azure"
  function_name = "header-auth"
  code_content  = "module.exports = async function (context, req) { return { status: 200 }; };"

  extra_config = {
    "azure_edge_rule_set" = "security-rules"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `function_name` (String, Required) Edge function name.
- `runtime` (String, Optional) Execution runtime (js-1.0, nodejs20).
- `code_content` (String, Required) Edge function source code.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
