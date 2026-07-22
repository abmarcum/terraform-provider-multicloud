---
# subcategory: "Compute"
page_title: "multicloud_serverless_function Resource - terraform-provider-multicloud"
description: |-
  Unified Serverless Event Function.
---

# multicloud_serverless_function (Resource)

Unified Serverless Event Function.

## Cloud Targets
- **AWS Target:** aws_lambda_function
- **GCP Target:** google_cloudfunctions2_function
- **Azure Target:** azurerm_linux_function_app

## How It Works

The `multicloud_serverless_function` resource provisions event-driven serverless code execution units across AWS Lambda, GCP Cloud Functions v2, and Azure Function Apps, managing runtimes, memory limits, and environment variables.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_serverless_function" "basic" {
  provider_type   = "aws"
  function_name   = "image-resizer"
  runtime         = "python3.11"
  handler         = "main.handler"
  memory_size_mb  = 512
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_serverless_function" "gcp_advanced" {
  provider_type   = "gcp"
  function_name   = "event-processor-gcp"
  runtime         = "nodejs20"
  handler         = "processEvent"
  memory_size_mb  = 1024
  timeout_seconds = 60

  extra_config = {
    "gcp_max_instance_count" = "50"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `function_name` (String, Required) Function name.
- `runtime` (String, Optional) Runtime environment.
- `handler` (String, Optional) Entrypoint handler.
- `memory_size_mb` (Int64, Optional) Memory in MB.
- `timeout_seconds` (Int64, Optional) Timeout in seconds.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
