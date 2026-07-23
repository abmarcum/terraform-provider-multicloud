---
# subcategory: "Network"
page_title: "multicloud_graphql_api Resource - terraform-provider-multicloud"
description: |-
  Unified Managed GraphQL API Endpoint.
---

# multicloud_graphql_api (Resource)

Unified Managed GraphQL API Endpoint.

## Cloud Targets
- **AWS Target:** aws_appsync_graphql_api
- **GCP Target:** google_apigee_environment
- **Azure Target:** azurerm_api_management_api

## How It Works

The `multicloud_graphql_api` resource provisions managed GraphQL gateways across AWS AppSync, GCP Apigee/Data Connect, and Azure API Management GraphQL APIs.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_graphql_api" "basic" {
  provider_type       = "aws"
  api_name            = "product-catalog-api"
  authentication_type = "API_KEY"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_graphql_api" "azure_advanced" {
  provider_type       = "azure"
  api_name            = "graphql-apim"
  authentication_type = "OIDC"

  extra_config = {
    "azure_path" = "graphql"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `api_name` (String, Required) GraphQL API name.
- `authentication_type` (String, Optional) Auth provider type.
- `schema_definition` (String, Optional) GraphQL SDL schema string.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
