---
# subcategory: "Database"
page_title: "multicloud_db_instance Resource - terraform-provider-multicloud"
description: |-
  Unified Relational Database Instance (SQL).
---

# multicloud_db_instance (Resource)

Unified Relational Database Instance (SQL).

## Cloud Targets
- **AWS Target:** aws_db_instance
- **GCP Target:** google_sql_database_instance
- **Azure Target:** azurerm_postgresql_server

## How It Works

The `multicloud_db_instance` resource provisions managed relational database instances across AWS RDS, GCP Cloud SQL, and Azure PostgreSQL/MySQL servers, configuring storage, engine versions, and Multi-AZ replication.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_db_instance" "basic" {
  provider_type        = "aws"
  instance_name        = "prod-db"
  engine               = "postgres"
  engine_version       = "15.2"
  size_tier            = "medium"
  storage_gb           = 100
  multi_az             = true
  backup_retention_days = 7
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_db_instance" "gcp_advanced" {
  provider_type        = "gcp"
  instance_name        = "cloudsql-postgres-gcp"
  engine               = "postgres"
  engine_version       = "14"
  size_tier            = "large"
  storage_gb           = 250
  multi_az             = true
  backup_retention_days = 14

  extra_config = {
    "gcp_database_flags" = "max_connections=200"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `instance_name` (String, Required) DB instance name.
- `engine` (String, Optional) 'postgres', 'mysql', or 'sqlserver'.
- `engine_version` (String, Optional) Version.
- `size_tier` (String, Optional) Size tier ('small', 'medium', 'large').
- `storage_gb` (Int64, Optional) Storage in GB.
- `multi_az` (Bool, Optional) High availability.
- `backup_retention_days` (Int64, Optional) Retention days.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
