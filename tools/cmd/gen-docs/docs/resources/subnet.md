---
# subcategory: "Network"
page_title: "multicloud_subnet Resource - terraform-provider-multicloud"
description: |-
  Unified Subnetwork.
---

# multicloud_subnet (Resource)

Unified Subnetwork.

## Cloud Targets
- **AWS Target:** aws_subnet
- **GCP Target:** google_compute_subnetwork
- **Azure Target:** azurerm_subnet

## How It Works

The `multicloud_subnet` resource provisions subnets within virtual networks across AWS Subnets, GCP Subnetworks, and Azure Subnets, associating IPv4 CIDR blocks and availability zone placement.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_subnet" "basic" {
  provider_type     = "gcp"
  subnet_name       = "public-subnet-1"
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-central1-a"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_subnet" "aws_advanced" {
  provider_type     = "aws"
  subnet_name       = "private-subnet-az2"
  cidr_block        = "10.0.2.0/24"
  availability_zone = "us-west-2b"

  extra_config = {
    "aws_map_public_ip_on_launch" = "false"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `subnet_name` (String, Required) Subnet identifier.
- `vpc_id` (String, Optional) Parent VPC ID.
- `cidr_block` (String, Optional) Subnet CIDR block.
- `availability_zone` (String, Optional) AZ placement.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
