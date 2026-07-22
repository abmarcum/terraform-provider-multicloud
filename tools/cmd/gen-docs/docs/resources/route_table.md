---
# subcategory: "Network"
page_title: "multicloud_route_table Resource - terraform-provider-multicloud"
description: |-
  Unified Routing Table.
---

# multicloud_route_table (Resource)

Unified Routing Table.

## Cloud Targets
- **AWS Target:** aws_route_table
- **GCP Target:** google_compute_route
- **Azure Target:** azurerm_route_table

## How It Works

The `multicloud_route_table` resource configures routing rules and gateway associations for virtual networks across AWS Route Tables, GCP Compute Routes, and Azure Route Tables.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_route_table" "basic" {
  provider_type = "aws"
  table_name    = "private-route-table"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_route_table" "advanced" {
  provider_type = "azure"
  table_name    = "spoke-vnet-routes"

  extra_config = {
    "azure_disable_bgp_route_propagation" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `table_name` (String, Required) Table identifier.
- `vpc_id` (String, Optional) Parent VPC ID.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
