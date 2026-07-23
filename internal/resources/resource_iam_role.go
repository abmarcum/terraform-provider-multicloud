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

var _ resource.Resource = &IAMRoleResource{}
var _ resource.ResourceWithImportState = &IAMRoleResource{}

type IAMRoleResource struct {
	clientManager interface{}
}

type IAMRoleModel struct {
	ID               types.String `tfsdk:"id"`
	RoleName         types.String `tfsdk:"role_name"`
	ProviderType     types.String `tfsdk:"provider_type"`
	Description      types.String `tfsdk:"description"`
	AssumeRolePolicy types.String `tfsdk:"assume_role_policy"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"` 
	Region       types.String `tfsdk:"region"` 
}

func NewIAMRoleResource() resource.Resource {
	return &IAMRoleResource{}
}

func (r *IAMRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iam_role"
}

func (r *IAMRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud IAM Role resource supporting AWS IAM Role, GCP Service Account, and Azure Managed Identity.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"role_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"assume_role_policy": schema.StringAttribute{
				Optional: true,
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *IAMRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *IAMRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IAMRoleModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "iam_role", plan.RoleName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.RoleName.IsUnknown() {
		if val, ok := res.Attributes["rolename"].(string); ok && val != "" {
			plan.RoleName = types.StringValue(val)
		} else {
			plan.RoleName = types.StringValue("default-rolename")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.Description.IsUnknown() {
		if val, ok := res.Attributes["description"].(string); ok && val != "" {
			plan.Description = types.StringValue(val)
		} else {
			plan.Description = types.StringValue("default-description")
		}
	}
	if plan.AssumeRolePolicy.IsUnknown() {
		if val, ok := res.Attributes["assumerolepolicy"].(string); ok && val != "" {
			plan.AssumeRolePolicy = types.StringValue(val)
		} else {
			plan.AssumeRolePolicy = types.StringValue("default-assumerolepolicy")
		}
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	if plan.Region.IsUnknown() {
		if val, ok := res.Attributes["region"].(string); ok && val != "" {
			plan.Region = types.StringValue(val)
		} else {
			plan.Region = types.StringValue("default-region")
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IAMRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IAMRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "us-central1"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	resName := state.RoleName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "iam_role", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IAMRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IAMRoleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !plan.ProviderType.IsNull() && plan.ProviderType.ValueString() != "" {
		pType = plan.ProviderType.ValueString()
	}
	reg := "us-central1"
	if !plan.Region.IsNull() && plan.Region.ValueString() != "" {
		reg = plan.Region.ValueString()
	}

	resName := plan.RoleName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "iam_role", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IAMRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IAMRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "iam_role", state.RoleName.ValueString(), reg)
}

func (r *IAMRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
