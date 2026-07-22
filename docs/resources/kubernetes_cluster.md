---
# subcategory: "Container"
page_title: "multicloud_kubernetes_cluster Resource - terraform-provider-multicloud"
description: |-
  Unified Managed Kubernetes Engine (EKS, GKE, AKS).
---

# multicloud_kubernetes_cluster (Resource)

Unified Managed Kubernetes Engine (EKS, GKE, AKS).

## Cloud Targets
- **AWS Target:** aws_eks_cluster
- **GCP Target:** google_container_cluster
- **Azure Target:** azurerm_kubernetes_cluster

## How It Works

The `multicloud_kubernetes_cluster` resource provisions managed Kubernetes control planes and node pools across AWS EKS, GCP GKE, and Azure AKS.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_kubernetes_cluster" "basic" {
  provider_type      = "aws"
  cluster_name       = "prod-eks-cluster"
  kubernetes_version = "1.30"
  node_count         = 6
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_kubernetes_cluster" "gcp_advanced" {
  provider_type      = "gcp"
  cluster_name       = "gke-autopilot-cluster"
  kubernetes_version = "1.29"

  extra_config = {
    "gcp_enable_autopilot" = "true"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `cluster_name` (String, Required) Cluster name.
- `kubernetes_version` (String, Optional) K8s version.
- `node_count` (Int64, Optional) Worker node count.
- `node_instance_type` (String, Optional) Node tier.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
