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
  name          = "company-shared-assets"
  region        = "us-west-2"
}

output "bucket_status" {
  value = data.multicloud_resource.existing_bucket.status
}
```

## Schema Attributes

### Required
- `name` (String) Name of the cloud resource to look up.
- `provider_type` (String) Target cloud provider (`'aws'`, `'gcp'`, or `'azure'`).
- `resource_type` (String) Unified resource type (e.g., `'storage_bucket'`, `'virtual_machine'`, `'db_instance'`).

### Optional
- `region` (String) Target cloud region.

### Read-Only
- `id` (String) Cloud-native identifier of the queried resource.
- `status` (String) Current lifecycle status of the cloud resource (`'ACTIVE'`, `'RUNNING'`, `'SUCCEEDED'`).
