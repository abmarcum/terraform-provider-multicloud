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

var _ resource.Resource = &ServerlessFunctionResource{}
var _ resource.ResourceWithImportState = &ServerlessFunctionResource{}

type ServerlessFunctionResource struct {
	clientManager interface{}
}

type ServerlessFunctionModel struct {
	ID                   types.String `tfsdk:"id"`
	FunctionName         types.String `tfsdk:"function_name"`
	ProviderType         types.String `tfsdk:"provider_type"`
	Runtime              types.String `tfsdk:"runtime"`
	Handler              types.String `tfsdk:"handler"`
	MemorySizeMB         types.Int64  `tfsdk:"memory_size_mb"`
	TimeoutSeconds       types.Int64  `tfsdk:"timeout_seconds"`
	Region               types.String `tfsdk:"region"`
	ExtraConfig          types.Map    `tfsdk:"extra_config"`
	EnvironmentVariables types.Map    `tfsdk:"environment_variables"`
}

func NewServerlessFunctionResource() resource.Resource {
	return &ServerlessFunctionResource{}
}

func (r *ServerlessFunctionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_serverless_function"
}

func (r *ServerlessFunctionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Serverless Function resource supporting AWS Lambda, GCP Cloud Functions, and Azure Functions.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"function_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"runtime": schema.StringAttribute{
				Optional: true,
			},
			"handler": schema.StringAttribute{
				Optional: true,
			},
			"memory_size_mb": schema.Int64Attribute{
				Optional: true,
			},
			"timeout_seconds": schema.Int64Attribute{
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
			"environment_variables": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *ServerlessFunctionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *ServerlessFunctionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServerlessFunctionModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "serverless_function", plan.FunctionName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.FunctionName.IsUnknown() {
		if val, ok := res.Attributes["functionname"].(string); ok && val != "" {
			plan.FunctionName = types.StringValue(val)
		} else {
			plan.FunctionName = types.StringValue("default-functionname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.Runtime.IsUnknown() {
		if val, ok := res.Attributes["runtime"].(string); ok && val != "" {
			plan.Runtime = types.StringValue(val)
		} else {
			plan.Runtime = types.StringValue("default-runtime")
		}
	}
	if plan.Handler.IsUnknown() {
		if val, ok := res.Attributes["handler"].(string); ok && val != "" {
			plan.Handler = types.StringValue(val)
		} else {
			plan.Handler = types.StringValue("default-handler")
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
	if plan.EnvironmentVariables.IsUnknown() {
		plan.EnvironmentVariables = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerlessFunctionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServerlessFunctionModel
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

	resName := state.FunctionName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "serverless_function", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServerlessFunctionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServerlessFunctionModel
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

	resName := plan.FunctionName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "serverless_function", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServerlessFunctionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServerlessFunctionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "serverless_function", state.FunctionName.ValueString(), reg)
}

func (r *ServerlessFunctionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
