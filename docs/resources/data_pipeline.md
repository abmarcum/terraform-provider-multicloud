---
# subcategory: "Analytics"
page_title: "multicloud_data_pipeline Resource - terraform-provider-multicloud"
description: |-
  Unified ETL & Batch Data Processing Pipeline.
---

# multicloud_data_pipeline (Resource)

Unified ETL & Batch Data Processing Pipeline.

## Cloud Targets
- **AWS Target:** aws_glue_crawler
- **GCP Target:** google_dataflow_job
- **Azure Target:** azurerm_data_factory_pipeline

## How It Works

The `multicloud_data_pipeline` resource provisions data extraction, transformation, and batch analytics processing jobs across AWS Glue/EMR, GCP Dataflow/Dataproc, and Azure Data Factory.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_data_pipeline" "basic" {
  provider_type = "aws"
  pipeline_name = "daily-etl-pipeline"
  engine        = "GLUE"
  max_workers   = 10
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_data_pipeline" "gcp_advanced" {
  provider_type = "gcp"
  pipeline_name = "streaming-dataflow"
  engine        = "FLINK"
  max_workers   = 20

  extra_config = {
    "gcp_temp_location" = "gs://my-bucket/tmp"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `pipeline_name` (String, Required) Pipeline job name.
- `engine` (String, Optional) Processing engine (SPARK, FLINK, GLUE).
- `max_workers` (Int64, Optional) Maximum worker node count.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
