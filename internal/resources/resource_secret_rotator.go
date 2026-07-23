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

var _ resource.Resource = &SecretRotatorResource{}
var _ resource.ResourceWithImportState = &SecretRotatorResource{}

type SecretRotatorResource struct {
	clientManager interface{}
}

type SecretRotatorModel struct {
	ID                   types.String `tfsdk:"id"`
	RotationName         types.String `tfsdk:"rotation_name"`
	ProviderType         types.String `tfsdk:"provider_type"`
	SecretID             types.String `tfsdk:"secret_id"`
	RotationIntervalDays types.Int64  `tfsdk:"rotation_interval_days"`
	RotationStatus       types.String `tfsdk:"rotation_status"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"` 
	Region       types.String `tfsdk:"region"` 
}

func NewSecretRotatorResource() resource.Resource {
	return &SecretRotatorResource{}
}

func (r *SecretRotatorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_rotator"
}

func (r *SecretRotatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Secret Rotator resource managing automated secret rotation across AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"rotation_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"secret_id": schema.StringAttribute{
				Required: true,
			},
			"rotation_interval_days": schema.Int64Attribute{
				Optional: true,
			},
			"rotation_status": schema.StringAttribute{
				Computed: true,
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

func (r *SecretRotatorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *SecretRotatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SecretRotatorModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "secret_rotator", plan.RotationName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.RotationName.IsUnknown() {
		if val, ok := res.Attributes["rotationname"].(string); ok && val != "" {
			plan.RotationName = types.StringValue(val)
		} else {
			plan.RotationName = types.StringValue("default-rotationname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.SecretID.IsUnknown() {
		if val, ok := res.Attributes["secretid"].(string); ok && val != "" {
			plan.SecretID = types.StringValue(val)
		} else {
			plan.SecretID = types.StringValue("default-secretid")
		}
	}
	if plan.RotationStatus.IsUnknown() {
		if val, ok := res.Attributes["rotationstatus"].(string); ok && val != "" {
			plan.RotationStatus = types.StringValue(val)
		} else {
			plan.RotationStatus = types.StringValue("default-rotationstatus")
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

func (r *SecretRotatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretRotatorModel
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

	resName := state.RotationName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "secret_rotator", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecretRotatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SecretRotatorModel
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

	resName := plan.RotationName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "secret_rotator", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecretRotatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SecretRotatorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "secret_rotator", state.RotationName.ValueString(), reg)
}

func (r *SecretRotatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
