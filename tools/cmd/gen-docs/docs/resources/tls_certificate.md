---
# subcategory: "Security"
page_title: "multicloud_tls_certificate Resource - terraform-provider-multicloud"
description: |-
  Unified Managed SSL/TLS Certificate supporting AWS ACM, GCP Certificate Manager, and Azure Key Vault Certificates.
---

# multicloud_tls_certificate (Resource)

Unified Managed SSL/TLS Certificate supporting AWS ACM, GCP Certificate Manager, and Azure Key Vault Certificates.

## Cloud Targets
- **AWS Target:** aws_acm_certificate
- **GCP Target:** google_certificate_manager_certificate
- **Azure Target:** azurerm_key_vault_certificate

## How It Works

The `multicloud_tls_certificate` resource provisions and auto-renews managed X.509 SSL/TLS certificates across AWS Certificate Manager (ACM), GCP Certificate Manager, and Azure Key Vault Certificates.

## Example Usage

### Basic Usage
```hcl
resource "multicloud_tls_certificate" "basic" {
  provider_type     = "aws"
  cert_name         = "prod-api-cert"
  domain_name       = "api.example.com"
  validation_method = "DNS"
}
```

### Advanced Usage with Cloud Escape Hatches (`extra_config`)
```hcl
resource "multicloud_tls_certificate" "gcp_advanced" {
  provider_type     = "gcp"
  cert_name         = "global-ingress-cert"
  domain_name       = "app.example.com"
  validation_method = "DNS"
  auto_renew        = true

  extra_config = {
    "gcp_scope" = "EDGE_CACHE"
  }
}
```

## Schema Attributes

### Required
- `provider_type` (String) Target cloud provider ('aws', 'gcp', or 'azure').

### Resource-Specific & Optional Attributes
- `cert_name` (String, Required) Identifier name of the managed TLS certificate.
- `domain_name` (String, Required) Primary FQDN domain name.
- `validation_method` (String, Optional) Domain validation method ('DNS' or 'EMAIL').
- `auto_renew` (Bool, Optional) Enable automated certificate renewal.
- `region` (String, Optional) Target placement region.
- `extra_config` (Map[String], Optional) Cloud-specific escape hatch key-value parameters passed through to upstream cloud SDKs.

### Read-Only
- `id` (String) State resource identifier (<cloud>/<region>/<name>).
