---
# subcategory: "Compute"
page_title: "multicloud_workflow Resource - terraform-provider-multicloud"
description: |-
  Unified Serverless Workflow Orchestration supporting AWS Step Functions, GCP Workflows, and Azure Logic Apps.
---

# multicloud_workflow (Resource)

Unified Serverless Workflow Orchestration supporting AWS Step Functions, GCP Workflows, and Azure Logic Apps.

## Cloud Targets
- **AWS Target:** aws_sfn_state_machine
- **GCP Target:** google_workflows_workflow
- **Azure Target:** azurerm_logic_app_workflow

## How It Works

The `multicloud_workflow` resource deploys stateful serverless workflow state machines across AWS Step Functions, GCP Workflows, and Azure Logic Apps.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_workflow" "basic" {
  provider_type = "aws"
  workflow_name = "order-fulfillment-sfn"
  definition    = "{\"StartAt\":\"ProcessOrder\",\"States\":{\"ProcessOrder\":{\"Type\":\"Pass\",\"End\":true}}}"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_workflow" "gcp_advanced" {
  provider_type = "gcp"
  workflow_name = "etl-orchestrator"
  workflow_type = "STANDARD"
  definition    = "main:\n  steps:\n    - init:\n        return: 'ok'"

  extra_config = {
    "gcp_call_log_level" = "LOG_ERRORS_ONLY"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `workflow_name` (String, Required) Name of the serverless workflow.
- `definition` (String, Required) Workflow state machine definition (JSON or YAML).
- `workflow_type` (String, Optional) Execution mode ('STANDARD' or 'EXPRESS').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
