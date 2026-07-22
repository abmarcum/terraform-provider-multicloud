---
# subcategory: "Database"
page_title: "multicloud_nosql_table Resource - terraform-provider-multicloud"
description: |-
  Unified NoSQL Document / Key-Value Database Table.
---

# multicloud_nosql_table (Resource)

Unified NoSQL Document / Key-Value Database Table.

## Cloud Targets
- **AWS Target:** aws_dynamodb_table
- **GCP Target:** google_firestore_database
- **Azure Target:** azurerm_cosmosdb_account

## How It Works

The `multicloud_nosql_table` resource provisions schemaless NoSQL tables across AWS DynamoDB, GCP Firestore, and Azure Cosmos DB.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_nosql_table" "basic" {
  provider_type = "aws"
  table_name    = "user-sessions"
  hash_key      = "session_id"
  billing_mode  = "PAY_PER_REQUEST"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_nosql_table" "azure_advanced" {
  provider_type = "azure"
  table_name    = "cosmos-events-azure"
  hash_key      = "event_id"

  extra_config = {
    "azure_consistency_level" = "Session"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `table_name` (String, Required) Table name.
- `hash_key` (String, Optional) Partition hash key.
- `range_key` (String, Optional) Sort range key.
- `billing_mode` (String, Optional) 'PAY_PER_REQUEST' or 'PROVISIONED'.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
