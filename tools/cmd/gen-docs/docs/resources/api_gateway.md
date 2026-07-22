---
# subcategory: "Network"
page_title: "multicloud_api_gateway Resource - terraform-provider-multicloud"
description: |-
  Unified API Gateway Endpoint Manager.
---

# multicloud_api_gateway (Resource)

Unified API Gateway Endpoint Manager.

## Cloud Targets
- **AWS Target:** aws_apigatewayv2_api
- **GCP Target:** google_api_gateway_gateway
- **Azure Target:** azurerm_api_management

## How It Works

The `multicloud_api_gateway` resource provisions API gateway endpoints across AWS API Gateway v2, GCP API Gateway, and Azure API Management (APIM), exposing public HTTPS endpoints.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_api_gateway" "basic" {
  provider_type = "aws"
  api_name      = "customer-api"
  protocol_type = "HTTP"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_api_gateway" "azure_advanced" {
  provider_type = "azure"
  api_name      = "apim-gateway-azure"
  protocol_type = "REST"

  extra_config = {
    "azure_publisher_email" = "admin@example.com"
    "azure_publisher_name"  = "Platform Team"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `api_name` (String, Required) API Gateway name.
- `protocol_type` (String, Optional) 'HTTP', 'REST', or 'WEBSOCKET'.
- `api_endpoint` (String, Read-Only) Computed HTTPS URL.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
