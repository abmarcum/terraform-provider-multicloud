---
# subcategory: "Compute"
page_title: "multicloud_container_app Resource - terraform-provider-multicloud"
description: |-
  Unified Serverless Container Application.
---

# multicloud_container_app (Resource)

Unified Serverless Container Application.

## Cloud Targets
- **AWS Target:** aws_apprunner_service
- **GCP Target:** google_cloud_run_v2_service
- **Azure Target:** azurerm_container_app

## How It Works

The `multicloud_container_app` resource provisions serverless container workloads across AWS App Runner / ECS Fargate, GCP Cloud Run v2, and Azure Container Apps.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_container_app" "basic" {
  provider_type = "gcp"
  app_name      = "api-service"
  image         = "gcr.io/my-project/api:latest"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_container_app" "aws_advanced" {
  provider_type = "aws"
  app_name      = "api-service-aws"
  image         = "public.ecr.aws/my-org/api:v1"

  extra_config = {
    "aws_auto_scaling_min_size" = "2"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `app_name` (String, Required) App service name.
- `image` (String, Optional) Container image URI.
- `cpu` (String, Optional) CPU spec.
- `memory` (String, Optional) Memory spec.
- `port` (Int64, Optional) Target port.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
