---
# subcategory: "Observability"
page_title: "multicloud_metric_alert Resource - terraform-provider-multicloud"
description: |-
  Unified Metric Threshold Alarm Rule.
---

# multicloud_metric_alert (Resource)

Unified Metric Threshold Alarm Rule.

## Cloud Targets
- **AWS Target:** aws_cloudwatch_metric_alarm
- **GCP Target:** google_monitoring_alert_policy
- **Azure Target:** azurerm_monitor_metric_alert

## How It Works

The `multicloud_metric_alert` resource provisions metric evaluation rules and threshold alarms across AWS CloudWatch Alarms, GCP Monitoring Alert Policies, and Azure Metric Alerts.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_metric_alert" "basic" {
  provider_type = "aws"
  alert_name    = "high-cpu-alarm"
  metric_name   = "CPUUtilization"
  threshold     = 85.0
  comparison    = "GreaterThanThreshold"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_metric_alert" "azure_advanced" {
  provider_type = "azure"
  alert_name    = "memory-alert"
  metric_name   = "Available Memory Bytes"
  threshold     = 1000000000
  comparison    = "LessThan"

  extra_config = {
    "azure_severity" = "1"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `alert_name` (String, Required) Alert rule name.
- `metric_name` (String, Optional) Target metric name.
- `threshold` (Float64, Optional) Trigger threshold.
- `comparison` (String, Optional) Comparison operator.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
