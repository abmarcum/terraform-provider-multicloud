package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProviderModel maps the provider configuration schema
type ProviderModel struct {
	AWS         *AWSConfigModel   `tfsdk:"aws"`
	GCP         *GCPConfigModel   `tfsdk:"gcp"`
	Azure       *AzureConfigModel `tfsdk:"azure"`
	DefaultTags types.Map         `tfsdk:"default_tags"`
}

type AWSConfigModel struct {
	Region    types.String `tfsdk:"region"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
	Profile   types.String `tfsdk:"profile"`
}

type GCPConfigModel struct {
	Project     types.String `tfsdk:"project"`
	Region      types.String `tfsdk:"region"`
	Credentials types.String `tfsdk:"credentials"`
}

type AzureConfigModel struct {
	SubscriptionID types.String `tfsdk:"subscription_id"`
	TenantID       types.String `tfsdk:"tenant_id"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	ResourceGroup  types.String `tfsdk:"resource_group"`
}
