---
# subcategory: "Database"
page_title: "multicloud_cache_cluster Resource - terraform-provider-multicloud"
description: |-
  Unified In-Memory Redis Cache Cluster.
---

# multicloud_cache_cluster (Resource)

Unified In-Memory Redis Cache Cluster.

## Cloud Targets
- **AWS Target:** aws_elasticache_cluster
- **GCP Target:** google_redis_instance
- **Azure Target:** azurerm_redis_cache

## How It Works

The `multicloud_cache_cluster` resource provisions managed in-memory caching clusters across AWS ElastiCache, GCP Memorystore, and Azure Cache for Redis.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_cache_cluster" "basic" {
  provider_type = "gcp"
  cluster_name  = "session-cache"
  engine        = "redis"
  num_nodes     = 2
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_cache_cluster" "aws_advanced" {
  provider_type = "aws"
  cluster_name  = "elasticache-redis-aws"
  engine        = "redis"
  node_type     = "cache.t4g.medium"
  num_nodes     = 3

  extra_config = {
    "aws_automatic_failover_enabled" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `cluster_name` (String, Required) Cache cluster name.
- `engine` (String, Optional) 'redis' or 'memcached'.
- `node_type` (String, Optional) Node tier.
- `num_nodes` (Int64, Optional) Node count.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
