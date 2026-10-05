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

var _ resource.Resource = &BlockVolumeResource{}
var _ resource.ResourceWithImportState = &BlockVolumeResource{}

type BlockVolumeResource struct {
	clientManager interface{}
}

type BlockVolumeModel struct {
	ID                types.String `tfsdk:"id"`
	VolumeName        types.String `tfsdk:"volume_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	SizeGB            types.Int64  `tfsdk:"size_gb"`
	VolumeType        types.String `tfsdk:"volume_type"`
	IOPS              types.Int64  `tfsdk:"iops"`
	EncryptionEnabled types.Bool   `tfsdk:"encryption_enabled"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewBlockVolumeResource() resource.Resource {
	return &BlockVolumeResource{}
}

func (r *BlockVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_volume"
}

func (r *BlockVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Block Storage Volume resource supporting AWS EBS Volumes, GCP Persistent Disks, and Azure Managed Disks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"volume_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the block storage volume.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"size_gb": schema.Int64Attribute{
				Required:    true,
				Description: "Volume capacity in gigabytes (GB).",
			},
			"volume_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Storage performance tier ('ssd', 'hdd', 'nvme').",
			},
			"iops": schema.Int64Attribute{
				Optional:    true,
				Description: "Provisioned IOPS for high-throughput volumes.",
			},
			"encryption_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable disk encryption at rest.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *BlockVolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *BlockVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BlockVolumeModel
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
		"size_gb": plan.SizeGB.ValueInt64(),
	}
	if !plan.EncryptionEnabled.IsNull() && !plan.EncryptionEnabled.IsUnknown() {
		extraAttrs["encryption_enabled"] = plan.EncryptionEnabled.ValueBool()
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "block_volume", plan.VolumeName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.VolumeType.IsUnknown() || plan.VolumeType.IsNull() {
		plan.VolumeType = types.StringValue("ssd")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BlockVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BlockVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "block_volume", state.VolumeName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BlockVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BlockVolumeModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "block_volume", plan.VolumeName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BlockVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BlockVolumeModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "block_volume", state.VolumeName.ValueString(), state.ExtraConfig, resp)
}

func (r *BlockVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
