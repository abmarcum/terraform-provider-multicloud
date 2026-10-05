package resources

import (
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &WAFPolicyResource{}
var _ resource.ResourceWithImportState = &WAFPolicyResource{}

type WAFPolicyResource struct {
	clientManager interface{}
}

type WAFPolicyModel struct {
	ID            types.String `tfsdk:"id"`
	PolicyName    types.String `tfsdk:"policy_name"`
	ProviderType  types.String `tfsdk:"provider_type"`
	DefaultAction types.String `tfsdk:"default_action"`
	Region        types.String `tfsdk:"region"`
	ExtraConfig   types.Map    `tfsdk:"extra_config"`
}

func NewWAFPolicyResource() resource.Resource {
	return &WAFPolicyResource{}
}

func (r *WAFPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_waf_policy"
}

func (r *WAFPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Web Application Firewall (WAF) Policy resource supporting AWS WAFv2 Web ACL, GCP Cloud Armor Security Policy, and Azure WAF Policy.",
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
			"default_action": schema.StringAttribute{
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

func (r *WAFPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *WAFPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WAFPolicyModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "waf_policy", plan.PolicyName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.PolicyName.IsUnknown() {
		if val, ok := res.Attributes["policyname"].(string); ok && val != "" {
			plan.PolicyName = types.StringValue(val)
		} else {
			plan.PolicyName = types.StringValue("default-policyname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.DefaultAction.IsUnknown() {
		if val, ok := res.Attributes["defaultaction"].(string); ok && val != "" {
			plan.DefaultAction = types.StringValue(val)
		} else {
			plan.DefaultAction = types.StringValue("default-defaultaction")
		}
	}
	if plan.Region.IsUnknown() {
		if val, ok := res.Attributes["region"].(string); ok && val != "" {
			plan.Region = types.StringValue(val)
		} else {
			plan.Region = types.StringValue("default-region")
		}
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WAFPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WAFPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "waf_policy", state.PolicyName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WAFPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WAFPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "waf_policy", plan.PolicyName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WAFPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WAFPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "waf_policy", state.PolicyName.ValueString(), state.ExtraConfig, resp)
}

func (r *WAFPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
