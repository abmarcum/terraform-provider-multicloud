package resources

import (
	"context"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/pricing"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &CloudResourceDataSource{}
var _ datasource.DataSource = &CostEstimateDataSource{}

// CloudResourceDataSource queries live or mock multi-cloud resource metadata across AWS, GCP, and Azure
type CloudResourceDataSource struct {
	clientManager interface{}
}

type CloudResourceDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ResourceType types.String `tfsdk:"resource_type"`
	ProviderType types.String `tfsdk:"provider_type"`
	Region       types.String `tfsdk:"region"`
	Status       types.String `tfsdk:"status"`
}

func NewCloudResourceDataSource() datasource.DataSource {
	return &CloudResourceDataSource{}
}

func (d *CloudResourceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource"
}

func (d *CloudResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Queries an existing unified multi-cloud resource across AWS, GCP, and Azure.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Cloud-native identifier of the queried resource.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the cloud resource to look up.",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "Unified resource type (e.g., 'storage_bucket', 'virtual_machine', 'db_instance').",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "Target cloud provider ('aws', 'gcp', or 'azure').",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Target cloud region.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current lifecycle status of the cloud resource.",
			},
		},
	}
}

func (d *CloudResourceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clientManager = req.ProviderData
}

func (d *CloudResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudResourceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType, reg := resolveProviderAndRegion(data.ProviderType, data.Region, d.clientManager)
	attrs := buildResourceExtraAttrs(d.clientManager, pType, types.MapNull(types.StringType), nil)

	res, err := adapters.ReadCloudResourceWithAttrs(ctx, pType, data.ResourceType.ValueString(), data.Name.ValueString(), reg, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Data Source Read Error", err.Error())
		return
	}

	data.ID = types.StringValue(res.ID)
	data.Region = types.StringValue(reg)
	data.Status = types.StringValue(res.Status)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// CostEstimateDataSource computes monthly cost estimates and ARM optimization recommendations
type CostEstimateDataSource struct{}

type CostEstimateDataSourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	ProviderType       types.String  `tfsdk:"provider_type"`
	ResourceType       types.String  `tfsdk:"resource_type"`
	SizeTier           types.String  `tfsdk:"size_tier"`
	MonthlyCostUSD     types.Float64 `tfsdk:"monthly_cost_usd"`
	SuggestedTier      types.String  `tfsdk:"suggested_tier"`
	EstimatedSavingUSD types.Float64 `tfsdk:"estimated_saving_usd"`
}

func NewCostEstimateDataSource() datasource.DataSource {
	return &CostEstimateDataSource{}
}

func (d *CostEstimateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cost_estimate"
}

func (d *CostEstimateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Calculates monthly cloud cost estimates and architecture optimization recommendations across AWS, GCP, and Azure.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the cost estimate query.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "Target cloud provider ('aws', 'gcp', or 'azure').",
			},
			"resource_type": schema.StringAttribute{
				Required:    true,
				Description: "Unified resource type (e.g., 'virtual_machine', 'db_instance', 'kubernetes_cluster').",
			},
			"size_tier": schema.StringAttribute{
				Optional:    true,
				Description: "Instance size tier ('small', 'medium', 'large').",
			},
			"monthly_cost_usd": schema.Float64Attribute{
				Computed:    true,
				Description: "Estimated monthly cost in USD.",
			},
			"suggested_tier": schema.StringAttribute{
				Computed:    true,
				Description: "Recommended ARM/Graviton/Tau/Ampere instance tier if applicable.",
			},
			"estimated_saving_usd": schema.Float64Attribute{
				Computed:    true,
				Description: "Estimated monthly savings in USD when switching to the suggested tier.",
			},
		},
	}
}

func (d *CostEstimateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CostEstimateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := strings.ToLower(data.ProviderType.ValueString())
	rType := strings.ToLower(data.ResourceType.ValueString())
	tier := ""
	if !data.SizeTier.IsNull() && !data.SizeTier.IsUnknown() {
		tier = data.SizeTier.ValueString()
	}

	cost := pricing.EstimateMonthlyCost(pType, rType, tier)
	data.ID = types.StringValue(pType + "/" + rType + "/" + tier)
	data.MonthlyCostUSD = types.Float64Value(cost)

	if rec := pricing.RecommendCostOptimizations(pType, rType, tier); rec != nil {
		data.SuggestedTier = types.StringValue(rec.SuggestedTier)
		data.EstimatedSavingUSD = types.Float64Value(rec.EstimatedSaving)
	} else {
		data.SuggestedTier = types.StringValue("")
		data.EstimatedSavingUSD = types.Float64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
