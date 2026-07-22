---
# subcategory: "Analytics"
page_title: "multicloud_data_warehouse Resource - terraform-provider-multicloud"
description: |-
  Unified Data Warehouse Analytical Cluster.
---

# multicloud_data_warehouse (Resource)

Unified Data Warehouse Analytical Cluster.

## Cloud Targets
- **AWS Target:** aws_redshift_cluster
- **GCP Target:** google_bigquery_dataset
- **Azure Target:** azurerm_synapse_workspace

## How It Works

The `multicloud_data_warehouse` resource provisions petabyte-scale analytical data warehouse clusters across AWS Redshift, GCP BigQuery Datasets, and Azure Synapse Workspaces.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_data_warehouse" "basic" {
  provider_type  = "aws"
  warehouse_name = "analytics-dw"
  node_type      = "ra3.4xlarge"
  num_nodes      = 4
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_data_warehouse" "gcp_advanced" {
  provider_type  = "gcp"
  warehouse_name = "bigquery_analytics_ds"

  extra_config = {
    "gcp_default_table_expiration_ms" = "3600000"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `warehouse_name` (String, Required) Analytical cluster name.
- `node_type` (String, Optional) Node spec.
- `num_nodes` (Int64, Optional) Node count.
- `database_name` (String, Optional) Initial database name.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
