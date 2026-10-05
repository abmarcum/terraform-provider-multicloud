---
# subcategory: "Analytics"
page_title: "multicloud_vector_index Resource - terraform-provider-multicloud"
description: |-
  Unified AI Vector Database Index supporting AWS OpenSearch Vector Engine, GCP Vertex AI Vector Search, and Azure AI Search.
---

# multicloud_vector_index (Resource)

Unified AI Vector Database Index supporting AWS OpenSearch Vector Engine, GCP Vertex AI Vector Search, and Azure AI Search.

## Cloud Targets
- **AWS Target:** aws_opensearchserverless_collection
- **GCP Target:** google_vertex_ai_index
- **Azure Target:** azurerm_search_service

## How It Works

The `multicloud_vector_index` resource provisions high-dimensional vector similarity search indexes for RAG and LLM embeddings across AWS OpenSearch Vector Engine, GCP Vertex AI Vector Search, and Azure AI Search.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_vector_index" "basic" {
  provider_type   = "gcp"
  index_name      = "rag-embeddings-index"
  dimensions      = 1536
  distance_metric = "COSINE"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_vector_index" "aws_advanced" {
  provider_type   = "aws"
  index_name      = "knowledge-base-vectors"
  dimensions      = 3072
  distance_metric = "DOT_PRODUCT"

  extra_config = {
    "aws_standby_replicas" = "ENABLED"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `index_name` (String, Required) Name of the vector search index.
- `dimensions` (Int64, Required) Vector embedding dimensionality (e.g., 768, 1536).
- `distance_metric` (String, Optional) Similarity metric ('COSINE', 'DOT_PRODUCT', 'EUCLIDEAN').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
