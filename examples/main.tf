terraform {
  required_providers {
    multicloud = {
      source  = "abmarcum/multicloud"
      version = "0.1.0"
    }
  }
}

provider "multicloud" {
  aws {
    region     = "us-west-2"
    access_key = var.aws_access_key
    secret_key = var.aws_secret_key
  }
  gcp {
    project     = "my-gcp-project"
    region      = "us-central1"
    credentials = var.gcp_credentials
  }
  azure {
    subscription_id = var.azure_subscription_id
    tenant_id       = var.azure_tenant_id
    client_id       = var.azure_client_id
    client_secret   = var.azure_client_secret
    resource_group  = "rg-multicloud"
  }

  default_tags = {
    environment = "production"
    managed_by  = "terraform"
    owner       = "platform-team"
  }
}

# 1. Unified Object Storage Bucket on AWS S3
resource "multicloud_storage_bucket" "aws_storage" {
  provider_type      = "aws"
  bucket_name        = "my-prod-data-bucket-aws"
  region             = "us-west-2"
  versioning_enabled = true
  encryption_enabled = true
}

# 2. Unified Object Storage Bucket on GCP Cloud Storage
resource "multicloud_storage_bucket" "gcp_storage" {
  provider_type      = "gcp"
  bucket_name        = "my-prod-data-bucket-gcp"
  region             = "us-central1"
  versioning_enabled = true
  encryption_enabled = true
}

# 3. Unified Object Storage Bucket on Azure Blob Storage
resource "multicloud_storage_bucket" "azure_storage" {
  provider_type      = "azure"
  bucket_name        = "myproddatabucketazure"
  region             = "eastus"
  versioning_enabled = true
  encryption_enabled = true
}

# 4. Unified Virtual Machine on AWS EC2
resource "multicloud_virtual_machine" "aws_vm" {
  provider_type = "aws"
  vm_name       = "app-server-aws"
  region        = "us-west-2"
  size_tier     = "medium"
}

# 5. Unified Virtual Machine on GCP Compute Engine
resource "multicloud_virtual_machine" "gcp_vm" {
  provider_type = "gcp"
  vm_name       = "app-server-gcp"
  region        = "us-central1"
  size_tier     = "medium"
}

# 6. Unified Secret in AWS Secrets Manager
resource "multicloud_secret" "aws_db_password" {
  provider_type = "aws"
  secret_name   = "db_master_password"
  secret_value  = "SuperSecretPassword123!"
}

# 7. Unified Secret in GCP Secret Manager
resource "multicloud_secret" "gcp_db_password" {
  provider_type = "gcp"
  secret_name   = "db_master_password"
  secret_value  = "SuperSecretPassword123!"
}
