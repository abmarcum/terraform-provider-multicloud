# Multi-Cloud Terraform Provider (`terraform-provider-multicloud`) User Guide

Welcome to the comprehensive user guide and operations runbook for `terraform-provider-multicloud`. This guide walks platform engineers and DevOps teams through authentication setups, HCL resource declarations, pre-apply security/cost intelligence, cloud-specific configuration pass-through (`extra_config`), migration from legacy cloud providers, and disaster recovery orchestration.

---

## 1. Authentication Setup

### 1.1 AWS Authentication (Ambient IRSA / Static Credentials)
The provider supports both ambient short-lived IAM credentials (recommended for EKS / EC2) and static credentials:

```hcl
provider "multicloud" {
  aws {
    region     = "us-west-2"
    access_key = var.aws_access_key # Optional if using ambient credentials / AWS IRSA
    secret_key = var.aws_secret_key # Optional
  }
}
```

### 1.2 GCP Authentication (Workload Identity / Service Account ADC)
Supports Application Default Credentials (ADC) or explicit JSON service account keys:

```hcl
provider "multicloud" {
  gcp {
    project     = "my-gcp-project-id"
    region      = "us-central1"
    credentials = file("${path.module}/service-account.json") # Optional if using ADC
  }
}
```

### 1.3 Azure Authentication (Entra ID / Client Secret)
Supports Azure Managed Service Identity (MSI) or service principal credentials:

```hcl
provider "multicloud" {
  azure {
    subscription_id = var.azure_subscription_id
    tenant_id       = var.azure_tenant_id
    client_id       = var.azure_client_id
    client_secret   = var.azure_client_secret
    resource_group  = "rg-multicloud-prod"
  }
}
```

### 1.4 Live Cloud API Execution vs Offline Mock Mode (`mock_mode`)
By default, `mock_mode` is **disabled (`false`)**. `terraform plan` and `terraform apply` execute live API calls directly against AWS SDK, GCP REST, and Azure ARM endpoints. If credentials are empty or invalid, authentic cloud diagnostics will be returned.

To enable offline simulation / CI/CD dry runs:

```hcl
provider "multicloud" {
  mock_mode = true # Enables offline synthetic provisioning
}
```

Or set the environment variable:
```bash
export MULTICLOUD_MOCK_MODE=true
```

---

## 2. Core HCL Provisioning & Cloud-Specific Pass-Through (`extra_config`)

Unified resources abstract 90%+ of standard cloud properties across AWS, GCP, and Azure. For cloud-specific escape-hatch parameters, use the optional **`extra_config` map block**:

```hcl
# AWS Storage Bucket with custom S3 Bucket Key Enabled
resource "multicloud_storage_bucket" "aws_data" {
  provider_type      = "aws"
  bucket_name        = "company-prod-storage-aws"
  region             = "us-west-2"
  versioning_enabled = true

  extra_config = {
    "aws_s3_bucket_key_enabled" = "true"
    "aws_force_destroy"         = "true"
  }
}

# GCP Storage Bucket with Custom Storage Class
resource "multicloud_storage_bucket" "gcp_data" {
  provider_type      = "gcp"
  bucket_name        = "company-prod-storage-gcp"
  region             = "us-central1"

  extra_config = {
    "gcp_storage_class" = "NEARLINE"
  }
}

# Intel Xeon Virtual Machine with Explicit Instance Type
resource "multicloud_virtual_machine" "app_server" {
  provider_type = "aws"
  vm_name       = "prod-app-server"
  region        = "us-west-2"
  size_tier     = "medium"      # Defaults to Intel Xeon Ice Lake (m6i.large)
  instance_type = "m6i.xlarge"  # Explicit Intel Xeon Platinum instance type SKU
}
```

---

## 3. Pre-Apply Security & Cost Intelligence

During `terraform plan`, the provider automatically executes pre-apply checks:

1. **CIS Benchmarks Auditor (`security_auditor.go`)**: Warns if a storage bucket is unencrypted or if a virtual machine is publicly exposed to `0.0.0.0/0`.
2. **Pre-Apply Monthly Cost Estimator (`cost_estimator.go`)**: Queries live public pricing APIs (Azure Retail REST API, AWS `AWS_PRICING_OFFER_URL`, GCP `GCP_BILLING_API_KEY`) with TTL caching and offline fallback resiliency to calculate monthly USD costs across all 70 resources.
3. **Cost Optimization Advisor (`cost_optimizer.go`)**: Suggests Arm64-based AWS Graviton3 (`t4g.*`), GCP Ampere Altra (`t2a-standard-*`), and Azure Cobalt 100 (`Standard_D*ps_v6`) instances for 20-30% cost savings.
4. **Secret Scanner (`secret_scanner.go`)**: Blocks hardcoded AWS secret keys, RSA private keys, and GCP service account JSON keys.
5. **Custom OPA Rego Policies (`opa_engine.go`)**: Point `export OPA_POLICY_PATH=/path/to/policies.rego` to evaluate custom `deny[msg]` and `violation[msg]` rules during `terraform plan`.
6. **OpenTelemetry Exporter (`telemetry.go`)**: Set `export OTEL_EXPORTER_OTLP_ENDPOINT=https://otel-collector:4318/v1/logs` to stream structured provision/read/update/delete events.

---

## 4. Data Sources (`multicloud_resource` & `multicloud_cost_estimate`)

Query existing cloud resources or pre-calculate FinOps estimates directly in HCL:

```hcl
data "multicloud_resource" "shared_bucket" {
  provider_type = "aws"
  resource_type = "storage_bucket"
  resource_name = "company-shared-assets"
  region        = "us-west-2"
}

data "multicloud_cost_estimate" "vm_estimate" {
  provider_type = "gcp"
  resource_type = "virtual_machine"
  size_tier     = "large"
}
```

---

## 5. Migrating Existing AWS, GCP, and Azure Infrastructure (`tf-migrate`)

Use the included **`tf-migrate` CLI tool** to convert existing legacy AWS (`aws_*`), GCP (`google_*`), and Azure (`azurerm_*`) `.tf` files into unified `multicloud_*` definitions. `tf-migrate` automatically extracts provider-specific properties into `extra_config` blocks and generates `.tfstate` migration scripts:

```bash
# Convert existing .tf files in ./legacy_tf and generate state migration script
go run ./tools/cmd/tf-migrate --input-dir ./legacy_tf --out-file multicloud_migrated.tf --migrate-state

# Preview converted HCL without modifying files
go run ./tools/cmd/tf-migrate --input-dir ./legacy_tf --dry-run
```

### Generated Migrated HCL Example:
```hcl
# Auto-generated by tf-migrate
resource "multicloud_storage_bucket" "prod_storage" {
  provider_type      = "aws"
  bucket_name        = "my-aws-bucket"
  versioning_enabled = true

  extra_config = {
    "aws_force_destroy" = "true"
  }
}
```
