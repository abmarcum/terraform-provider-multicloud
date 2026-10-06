---
# subcategory: "Compute"
page_title: "multicloud_batch_compute Resource - terraform-provider-multicloud"
description: |-
  Unified Batch Compute Environment supporting AWS Batch, GCP Cloud Batch, and Azure Batch Pools.
---

# multicloud_batch_compute (Resource)

Unified Batch Compute Environment supporting AWS Batch, GCP Cloud Batch, and Azure Batch Pools.

## Cloud Targets
- **AWS Target:** aws_batch_compute_environment
- **GCP Target:** google_batch_job
- **Azure Target:** azurerm_batch_pool

## How It Works

The `multicloud_batch_compute` resource provisions managed high-performance computing (HPC) and batch job execution pools across AWS Batch, GCP Cloud Batch, and Azure Batch.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_batch_compute" "basic" {
  provider_type    = "aws"
  environment_name = "genomics-batch-pool"
  max_vcpus        = 256
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_batch_compute" "azure_advanced" {
  provider_type    = "azure"
  environment_name = "rendering-spot-pool"
  max_vcpus        = 512
  compute_type     = "SPOT"

  extra_config = {
    "azure_inter_node_communication" = "Enabled"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `environment_name` (String, Required) Name of the batch compute environment.
- `max_vcpus` (Int64, Required) Maximum vCPU capacity.
- `compute_type` (String, Optional) Provisioning model ('EC2', 'FARGATE', 'SPOT').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
