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

var _ resource.Resource = &AppConfigResource{}
var _ resource.ResourceWithImportState = &AppConfigResource{}

type AppConfigResource struct {
	clientManager interface{}
}

type AppConfigModel struct {
	ID           types.String `tfsdk:"id"`
	ConfigName   types.String `tfsdk:"config_name"`
	ProviderType types.String `tfsdk:"provider_type"`
	ConfigKey    types.String `tfsdk:"config_key"`
	ConfigValue  types.String `tfsdk:"config_value"`
	Region       types.String `tfsdk:"region"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"`
}

func NewAppConfigResource() resource.Resource {
	return &AppConfigResource{}
}

func (r *AppConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_config"
}

func (r *AppConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Key-Value Application Configuration Store resource supporting AWS SSM Parameter, GCP Runtime Config, and Azure App Configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"config_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"config_key": schema.StringAttribute{
				Optional: true,
			},
			"config_value": schema.StringAttribute{
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

func (r *AppConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *AppConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppConfigModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "app_config", plan.ConfigName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.ConfigName.IsUnknown() {
		if val, ok := res.Attributes["configname"].(string); ok && val != "" {
			plan.ConfigName = types.StringValue(val)
		} else {
			plan.ConfigName = types.StringValue("default-configname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.ConfigKey.IsUnknown() {
		if val, ok := res.Attributes["configkey"].(string); ok && val != "" {
			plan.ConfigKey = types.StringValue(val)
		} else {
			plan.ConfigKey = types.StringValue("default-configkey")
		}
	}
	if plan.ConfigValue.IsUnknown() {
		if val, ok := res.Attributes["configvalue"].(string); ok && val != "" {
			plan.ConfigValue = types.StringValue(val)
		} else {
			plan.ConfigValue = types.StringValue("default-configvalue")
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

func (r *AppConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppConfigModel
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

	resName := state.ConfigName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "app_config", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppConfigModel
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

	resName := plan.ConfigName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "app_config", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppConfigModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "app_config", state.ConfigName.ValueString(), reg)
}

func (r *AppConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
