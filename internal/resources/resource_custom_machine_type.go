package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CustomMachineTypeResource{}
var _ resource.ResourceWithImportState = &CustomMachineTypeResource{}

type CustomMachineTypeResource struct {
	clientManager interface{}
}

type CustomMachineTypeModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ProviderType     types.String `tfsdk:"provider_type"`
	Region           types.String `tfsdk:"region"`
	VCPUs            types.Int64  `tfsdk:"vcpus"`
	MemoryMB         types.Int64  `tfsdk:"memory_mb"`
	MachineTypeSpec  types.String `tfsdk:"machine_type_spec"`
	ExtraConfig      types.Map    `tfsdk:"extra_config"`
}

func NewCustomMachineTypeResource() resource.Resource {
	return &CustomMachineTypeResource{}
}

func (r *CustomMachineTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_machine_type"
}

func (r *CustomMachineTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Custom Machine Type resource supporting GCP Custom Machine Types (custom-vCPU-RAM), AWS EC2 Custom Launch Templates, and Azure Custom VM sizes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
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
			"vcpus": schema.Int64Attribute{
				Required:    true,
				Description: "Number of custom vCPUs (e.g., 2, 4, 8, 16).",
			},
			"memory_mb": schema.Int64Attribute{
				Required:    true,
				Description: "Custom memory allocation in megabytes (e.g., 4096, 8192, 16384).",
			},
			"machine_type_spec": schema.StringAttribute{
				Computed:    true,
				Description: "Formatted custom machine type specification string (e.g. 'custom-4-16384').",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *CustomMachineTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *CustomMachineTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CustomMachineTypeModel
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
	extraAttrs["vcpus"] = plan.VCPUs.ValueInt64()
	extraAttrs["memory_mb"] = plan.MemoryMB.ValueInt64()

	res, err := adapters.CreateCloudResource(ctx, providerType, "custom_machine_type", plan.Name.ValueString(), reg, extraAttrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Name.IsUnknown() {
		plan.Name = types.StringValue("default-custom-machine")
	}
	if plan.ProviderType.IsUnknown() {
		plan.ProviderType = types.StringValue(providerType)
	}
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	plan.MachineTypeSpec = types.StringValue(fmt.Sprintf("custom-%d-%d", plan.VCPUs.ValueInt64(), plan.MemoryMB.ValueInt64()))
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomMachineTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CustomMachineTypeModel
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

	resName := state.Name.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "custom_machine_type", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CustomMachineTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CustomMachineTypeModel
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

	resName := plan.Name.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "custom_machine_type", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CustomMachineTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CustomMachineTypeModel
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

	_ = adapters.DeleteCloudResource(ctx, pType, "custom_machine_type", state.Name.ValueString(), reg)
}

func (r *CustomMachineTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
