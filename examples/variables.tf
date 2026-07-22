variable "aws_access_key" {
  type      = string
  default   = "mock_aws_access_key"
  sensitive = true
}

variable "aws_secret_key" {
  type      = string
  default   = "mock_aws_secret_key"
  sensitive = true
}

variable "gcp_credentials" {
  type      = string
  default   = "{}"
  sensitive = true
}

variable "azure_subscription_id" {
  type    = string
  default = "00000000-0000-0000-0000-000000000000"
}

variable "azure_tenant_id" {
  type      = string
  default   = "00000000-0000-0000-0000-000000000000"
  sensitive = true
}

variable "azure_client_id" {
  type      = string
  default   = "00000000-0000-0000-0000-000000000000"
  sensitive = true
}

variable "azure_client_secret" {
  type      = string
  default   = "mock_azure_client_secret"
  sensitive = true
}
