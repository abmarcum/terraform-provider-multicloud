---
# subcategory: "Analytics"
page_title: "multicloud_search_index Resource - terraform-provider-multicloud"
description: |-
  Unified Search Engine Index Domain.
---

# multicloud_search_index (Resource)

Unified Search Engine Index Domain.

## Cloud Targets
- **AWS Target:** aws_opensearch_domain
- **GCP Target:** google_discovery_engine_search_engine
- **Azure Target:** azurerm_search_service

## How It Works

The `multicloud_search_index` resource provisions managed search engine clusters across AWS OpenSearch, GCP Vertex AI Search, and Azure AI Search Services.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_search_index" "basic" {
  provider_type  = "aws"
  index_name     = "app-logs-index"
  instance_count = 3
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_search_index" "advanced" {
  provider_type  = "aws"
  index_name     = "es-prod-cluster"
  instance_type  = "t3.medium.search"
  instance_count = 5

  extra_config = {
    "aws_dedicated_master_enabled" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `index_name` (String, Required) Search index domain name.
- `instance_type` (String, Optional) Instance spec.
- `instance_count` (Int64, Optional) Node count.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
