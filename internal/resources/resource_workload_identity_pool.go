package resources

import (
	"context"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &WorkloadIdentityPoolResource{}
var _ resource.ResourceWithImportState = &WorkloadIdentityPoolResource{}

type WorkloadIdentityPoolResource struct {
	clientManager interface{}
}

type WorkloadIdentityPoolModel struct {
	ID               types.String `tfsdk:"id"`
	PoolName         types.String `tfsdk:"pool_name"`
	ProviderType     types.String `tfsdk:"provider_type"`
	Region           types.String `tfsdk:"region"`
	IssuerURL        types.String `tfsdk:"issuer_url"`
	AllowedAudiences types.List   `tfsdk:"allowed_audiences"`
	Description      types.String `tfsdk:"description"`
	Disabled         types.Bool   `tfsdk:"disabled"`
	ExtraConfig      types.Map    `tfsdk:"extra_config"`
}

func NewWorkloadIdentityPoolResource() resource.Resource {
	return &WorkloadIdentityPoolResource{}
}

func (r *WorkloadIdentityPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workload_identity_pool"
}

func (r *WorkloadIdentityPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Workload Identity Pool resource supporting GCP Workload Identity Pools, AWS IAM OIDC Providers, and Azure Entra ID Federated Credentials.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"pool_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"issuer_url": schema.StringAttribute{
				Required:    true,
				Description: "OIDC Issuer URL (e.g. 'https://token.actions.githubusercontent.com').",
			},
			"allowed_audiences": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Allowed OIDC audience client IDs.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description of the workload identity pool.",
			},
			"disabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether the workload identity pool is disabled (default false).",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *WorkloadIdentityPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *WorkloadIdentityPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkloadIdentityPoolModel
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

	extraAttrs := make(map[string]interface{})
	extraAttrs["issuer_url"] = plan.IssuerURL.ValueString()

	res, err := adapters.CreateCloudResource(ctx, providerType, "workload_identity_pool", plan.PoolName.ValueString(), reg, extraAttrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.PoolName.IsUnknown() {
		plan.PoolName = types.StringValue("default-identity-pool")
	}
	if plan.ProviderType.IsUnknown() {
		plan.ProviderType = types.StringValue(providerType)
	}
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("global")
	}
	if plan.Disabled.IsUnknown() || plan.Disabled.IsNull() {
		plan.Disabled = types.BoolValue(false)
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkloadIdentityPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkloadIdentityPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "global"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	resName := state.PoolName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "workload_identity_pool", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WorkloadIdentityPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkloadIdentityPoolModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !plan.ProviderType.IsNull() && plan.ProviderType.ValueString() != "" {
		pType = plan.ProviderType.ValueString()
	}
	reg := "global"
	if !plan.Region.IsNull() && plan.Region.ValueString() != "" {
		reg = plan.Region.ValueString()
	}

	resName := plan.PoolName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "workload_identity_pool", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WorkloadIdentityPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkloadIdentityPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "global"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	_ = adapters.DeleteCloudResource(ctx, pType, "workload_identity_pool", state.PoolName.ValueString(), reg)
}

func (r *WorkloadIdentityPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
