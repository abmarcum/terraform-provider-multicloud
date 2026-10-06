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

var _ resource.Resource = &BatchComputeResource{}
var _ resource.ResourceWithImportState = &BatchComputeResource{}

type BatchComputeResource struct {
	clientManager interface{}
}

type BatchComputeModel struct {
	ID              types.String `tfsdk:"id"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	ProviderType    types.String `tfsdk:"provider_type"`
	Region          types.String `tfsdk:"region"`
	MaxVCPUs        types.Int64  `tfsdk:"max_vcpus"`
	ComputeType     types.String `tfsdk:"compute_type"`
	ExtraConfig     types.Map    `tfsdk:"extra_config"`
}

func NewBatchComputeResource() resource.Resource {
	return &BatchComputeResource{}
}

func (r *BatchComputeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_batch_compute"
}

func (r *BatchComputeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Batch Compute Environment resource supporting AWS Batch, GCP Cloud Batch, and Azure Batch Pools.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"environment_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the batch compute environment or pool.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"max_vcpus": schema.Int64Attribute{
				Required:    true,
				Description: "Maximum vCPU capacity for the batch compute environment.",
			},
			"compute_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Provisioning model ('EC2', 'FARGATE', 'SPOT').",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *BatchComputeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *BatchComputeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BatchComputeModel
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

	extraAttrs := map[string]interface{}{
		"max_vcpus": plan.MaxVCPUs.ValueInt64(),
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "batch_compute", plan.EnvironmentName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.ComputeType.IsUnknown() || plan.ComputeType.IsNull() {
		plan.ComputeType = types.StringValue("EC2")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BatchComputeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BatchComputeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "batch_compute", state.EnvironmentName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BatchComputeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BatchComputeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "batch_compute", plan.EnvironmentName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BatchComputeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BatchComputeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "batch_compute", state.EnvironmentName.ValueString(), state.ExtraConfig, resp)
}

func (r *BatchComputeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
