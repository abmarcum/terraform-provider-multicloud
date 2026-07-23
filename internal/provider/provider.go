package provider

import (
	"context"

	"github.com/abmarcum/multi-cloud-provider/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &MulticloudProvider{}

type MulticloudProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MulticloudProvider{
			version: version,
		}
	}
}

func (p *MulticloudProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "multicloud"
	resp.Version = p.version
}

func (p *MulticloudProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Terraform Provider for AWS, GCP, and Azure under a single unified schema.",
		Blocks: map[string]schema.Block{
			"aws": schema.SingleNestedBlock{
				Description: "AWS Cloud provider authentication & region settings.",
				Attributes: map[string]schema.Attribute{
					"region":     schema.StringAttribute{Optional: true, Description: "AWS Region (e.g. us-west-2)."},
					"access_key": schema.StringAttribute{Optional: true, Sensitive: true, Description: "AWS Access Key ID."},
					"secret_key": schema.StringAttribute{Optional: true, Sensitive: true, Description: "AWS Secret Access Key."},
					"profile":    schema.StringAttribute{Optional: true, Description: "AWS Profile name from ~/.aws/credentials."},
				},
			},
			"gcp": schema.SingleNestedBlock{
				Description: "GCP Cloud provider authentication & project settings.",
				Attributes: map[string]schema.Attribute{
					"project":     schema.StringAttribute{Optional: true, Description: "GCP Project ID."},
					"region":      schema.StringAttribute{Optional: true, Description: "GCP Region (e.g. us-central1)."},
					"credentials": schema.StringAttribute{Optional: true, Sensitive: true, Description: "GCP Service Account JSON string or file path."},
				},
			},
			"azure": schema.SingleNestedBlock{
				Description: "Azure Cloud provider authentication settings.",
				Attributes: map[string]schema.Attribute{
					"subscription_id": schema.StringAttribute{Optional: true, Description: "Azure Subscription ID."},
					"tenant_id":       schema.StringAttribute{Optional: true, Sensitive: true, Description: "Azure Tenant ID."},
					"client_id":       schema.StringAttribute{Optional: true, Sensitive: true, Description: "Azure Client ID."},
					"client_secret":   schema.StringAttribute{Optional: true, Sensitive: true, Description: "Azure Client Secret."},
					"resource_group":  schema.StringAttribute{Optional: true, Description: "Default Azure Resource Group name."},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"default_tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Default tags to apply across all provisioned resources (AWS Tags, GCP Labels, Azure Tags).",
			},
		},
	}
}

func (p *MulticloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientManager, err := NewClientManager(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to initialize cloud clients", err.Error())
		return
	}

	resp.DataSourceData = clientManager
	resp.ResourceData = clientManager
}

func (p *MulticloudProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewStorageBucketResource,
		resources.NewVirtualMachineResource,
		resources.NewVirtualNetworkResource,
		resources.NewSubnetResource,
		resources.NewSecurityGroupResource,
		resources.NewStaticIPResource,
		resources.NewNATGatewayResource,
		resources.NewRouteTableResource,
		resources.NewLoadBalancerResource,
		resources.NewDBInstanceResource,
		resources.NewNoSQLTableResource,
		resources.NewKubernetesClusterResource,
		resources.NewContainerRegistryResource,
		resources.NewServerlessFunctionResource,
		resources.NewSecretResource,
		resources.NewKMSKeyResource,
		resources.NewIAMRoleResource,
		resources.NewDNSZoneResource,
		resources.NewPubSubTopicResource,
		resources.NewMessageQueueResource,
		resources.NewFailoverPolicyResource,
		resources.NewEventBridgeResource,
		resources.NewCDNDistributionResource,
		resources.NewCacheClusterResource,
		resources.NewAPIGatewayResource,
		resources.NewDataWarehouseResource,
		resources.NewVPNGatewayResource,
		resources.NewSearchIndexResource,
		resources.NewAutoScalingGroupResource,
		resources.NewMonitoringDashboardResource,
		resources.NewDataSyncResource,
		resources.NewIdentityFederationResource,
		resources.NewSecretRotatorResource,
		resources.NewContainerAppResource,
		resources.NewBastionHostResource,
		resources.NewWAFPolicyResource,
		resources.NewVPCPeeringResource,
		resources.NewAppConfigResource,
		resources.NewAIEndpointResource,
		resources.NewStreamingClusterResource,
		resources.NewMetricAlertResource,
		resources.NewLogWorkspaceResource,
		resources.NewGraphQLAPIResource,
	}
}

func (p *MulticloudProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
