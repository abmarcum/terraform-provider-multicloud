---
# subcategory: "Network"
page_title: "multicloud_load_balancer Resource - terraform-provider-multicloud"
description: |-
  Unified Network / Application Load Balancer.
---

# multicloud_load_balancer (Resource)

Unified Network / Application Load Balancer.

## Cloud Targets
- **AWS Target:** aws_lb
- **GCP Target:** google_compute_forwarding_rule
- **Azure Target:** azurerm_lb

## How It Works

The `multicloud_load_balancer` resource provisions layer 4 (Network) and layer 7 (Application) load balancers across AWS ALB/NLB, GCP Forwarding Rules, and Azure Load Balancers.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_load_balancer" "basic" {
  provider_type = "gcp"
  balancer_name = "ingress-lb"
  balancer_type = "application"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_load_balancer" "aws_advanced" {
  provider_type = "aws"
  balancer_name = "alb-external-aws"
  balancer_type = "application"

  extra_config = {
    "aws_idle_timeout"             = "60"
    "aws_enable_deletion_protection" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `balancer_name` (String, Required) Load balancer name.
- `balancer_type` (String, Optional) 'application' or 'network'.
- `subnet_ids` (List[String], Optional) Target Subnet IDs.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
