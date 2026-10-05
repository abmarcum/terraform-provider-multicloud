---
# subcategory: "Data Sources"
page_title: "multicloud_resource Data Source - terraform-provider-multicloud"
description: |-
  Reads live or mock state and attributes for any existing multi-cloud resource across AWS, GCP, and Azure.
---

# multicloud_resource (Data Source)

Reads upstream state and metadata attributes for an existing cloud resource across **AWS**, **GCP**, or **Azure** using [`ReadCloudResourceWithAttrs`](../../internal/cloud/adapters/adapter_interface.go).

## Example Usage

```hcl
data "multicloud_resource" "existing_bucket" {
  provider_type = "aws"
  resource_type = "storage_bucket"
  resource_name = "company-shared-assets"
  region        = "us-west-2"
}

output "bucket_status" {
  value = data.multicloud_resource.existing_bucket.status
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider (`'aws'`, `'gcp'`, or `'azure'`).
- `resource_type` (String) Unified resource type (e.g., `'storage_bucket'`, `'virtual_machine'`, `'db_instance'`).
- `resource_name` (String) Upstream identifier or name of the cloud resource.

### Optional
- `region` (String) Target cloud region.

### Read-Only
- `id` (String) Resolved cloud resource identifier.
- `status` (String) Upstream cloud resource status (`'ACTIVE'`, `'RUNNING'`, `'SUCCEEDED'`).
- `attributes` (Map[String]) Key-value attributes returned by the cloud provider.
