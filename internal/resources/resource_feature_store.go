package resources

import (
	"context"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FeatureStoreResource{}
var _ resource.ResourceWithImportState = &FeatureStoreResource{}

type FeatureStoreResource struct {
	clientManager interface{}
}

type FeatureStoreModel struct {
	ID                 types.String `tfsdk:"id"`
	StoreName          types.String `tfsdk:"store_name"`
	ProviderType       types.String `tfsdk:"provider_type"`
	OnlineStoreEnabled types.Bool   `tfsdk:"online_store_enabled"`
	Region             types.String `tfsdk:"region"`
	ExtraConfig        types.Map    `tfsdk:"extra_config"`
}

func NewFeatureStoreResource() resource.Resource {
	return &FeatureStoreResource{}
}

func (r *FeatureStoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_feature_store"
}

func (r *FeatureStoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud MLOps Feature Store resource supporting AWS SageMaker Feature Group, GCP Vertex AI Featurestore, and Azure ML Workspace.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"store_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"online_store_enabled": schema.BoolAttribute{
				Optional: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *FeatureStoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *FeatureStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FeatureStoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	providerType := strings.ToLower(plan.ProviderType.ValueString())
	reg := ""
	if !plan.Region.IsNull() && !plan.Region.IsUnknown() {
		reg = plan.Region.ValueString()
	} else {
		plan.Region = types.StringNull()
	}
	res, err := adapters.CreateCloudResource(ctx, providerType, "feature_store", plan.StoreName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FeatureStoreModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "feature_store", state.StoreName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FeatureStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FeatureStoreModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "feature_store", plan.StoreName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FeatureStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FeatureStoreModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "feature_store", state.StoreName.ValueString(), state.ExtraConfig, resp)
}

func (r *FeatureStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
