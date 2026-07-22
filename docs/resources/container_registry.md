---
# subcategory: "Container"
page_title: "multicloud_container_registry Resource - terraform-provider-multicloud"
description: |-
  Unified Private Container Image Registry.
---

# multicloud_container_registry (Resource)

Unified Private Container Image Registry.

## Cloud Targets
- **AWS Target:** aws_ecr_repository
- **GCP Target:** google_artifact_registry_repository
- **Azure Target:** azurerm_container_registry

## How It Works

The `multicloud_container_registry` resource provisions private OCI container registries across AWS ECR, GCP Artifact Registry, and Azure Container Registry (ACR).

## Example Usage

### Basic Usage
```hcl
resource "multicloud_container_registry" "basic" {
  provider_type        = "aws"
  registry_name        = "app-images"
  image_tag_mutability = "IMMUTABLE"
  scan_on_push         = true
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_container_registry" "azure_advanced" {
  provider_type        = "azure"
  registry_name        = "acrprodregistry"
  image_tag_mutability = "IMMUTABLE"

  extra_config = {
    "azure_admin_enabled" = "false"
    "azure_sku"           = "Premium"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `registry_name` (String, Required) Registry name.
- `image_tag_mutability` (String, Optional) 'MUTABLE' or 'IMMUTABLE'.
- `scan_on_push` (Bool, Optional) Scan on image push.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
