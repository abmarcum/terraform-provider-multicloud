---
# subcategory: "Networking"
page_title: "multicloud_service_mesh Resource - terraform-provider-multicloud"
description: |-
  Unified Service Mesh Control Plane supporting AWS App Mesh, GCP Cloud Service Mesh, and Azure Kubernetes Fleet Service Mesh.
---

# multicloud_service_mesh (Resource)

Unified Service Mesh Control Plane supporting AWS App Mesh, GCP Cloud Service Mesh, and Azure Kubernetes Fleet Service Mesh.

## Cloud Targets
- **AWS Target:** aws_appmesh_mesh
- **GCP Target:** google_network_services_mesh
- **Azure Target:** azurerm_kubernetes_fleet_manager

## How It Works

The `multicloud_service_mesh` resource provisions zero-trust L7 service mesh control planes with mTLS and egress filtering across AWS App Mesh, GCP Cloud Service Mesh, and Azure Fleet Service Mesh.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_service_mesh" "basic" {
  provider_type = "aws"
  mesh_name     = "prod-microservices-mesh"
  mtls_mode     = "STRICT"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_service_mesh" "gcp_advanced" {
  provider_type = "gcp"
  mesh_name     = "global-zero-trust-mesh"
  mtls_mode     = "STRICT"
  egress_filter = "DROP_ALL"

  extra_config = {
    "gcp_interception_port" = "15001"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `mesh_name` (String, Required) Name of the service mesh control plane.
- `mtls_mode` (String, Optional) Mutual TLS policy ('STRICT' or 'PERMISSIVE').
- `egress_filter` (String, Optional) Outbound traffic filter ('ALLOW_ALL' or 'DROP_ALL').
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
