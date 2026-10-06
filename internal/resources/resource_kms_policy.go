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

var _ resource.Resource = &KMSPolicyResource{}
var _ resource.ResourceWithImportState = &KMSPolicyResource{}

type KMSPolicyResource struct {
	clientManager interface{}
}

type KMSPolicyModel struct {
	ID           types.String `tfsdk:"id"`
	PolicyName   types.String `tfsdk:"policy_name"`
	ProviderType types.String `tfsdk:"provider_type"`
	KeyID        types.String `tfsdk:"key_id"`
	PolicyJSON   types.String `tfsdk:"policy_json"`
	Region       types.String `tfsdk:"region"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"`
}

func NewKMSPolicyResource() resource.Resource {
	return &KMSPolicyResource{}
}

func (r *KMSPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kms_policy"
}

func (r *KMSPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud KMS Policy resource supporting AWS KMS Key Policy / Grants, GCP KMS IAM Bindings, and Azure Key Vault Access Policies.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"policy_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"key_id": schema.StringAttribute{
				Required: true,
			},
			"policy_json": schema.StringAttribute{
				Required: true,
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

func (r *KMSPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *KMSPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan KMSPolicyModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "kms_policy", plan.PolicyName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
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

func (r *KMSPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state KMSPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "kms_policy", state.PolicyName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *KMSPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan KMSPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "kms_policy", plan.PolicyName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *KMSPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state KMSPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "kms_policy", state.PolicyName.ValueString(), state.ExtraConfig, resp)
}

func (r *KMSPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
