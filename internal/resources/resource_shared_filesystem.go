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

var _ resource.Resource = &SharedFilesystemResource{}
var _ resource.ResourceWithImportState = &SharedFilesystemResource{}

type SharedFilesystemResource struct {
	clientManager interface{}
}

type SharedFilesystemModel struct {
	ID                types.String `tfsdk:"id"`
	FilesystemName    types.String `tfsdk:"filesystem_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	Protocol          types.String `tfsdk:"protocol"`
	PerformanceMode   types.String `tfsdk:"performance_mode"`
	EncryptionEnabled types.Bool   `tfsdk:"encryption_enabled"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewSharedFilesystemResource() resource.Resource {
	return &SharedFilesystemResource{}
}

func (r *SharedFilesystemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shared_filesystem"
}

func (r *SharedFilesystemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Shared Network Filesystem resource supporting AWS EFS, GCP Cloud Filestore, and Azure Files.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"filesystem_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the shared network filesystem.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"protocol": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Network file sharing protocol ('NFSv4' or 'SMB').",
			},
			"performance_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Throughput mode ('generalPurpose' or 'maxIO').",
			},
			"encryption_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable encryption at rest for shared file storage.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *SharedFilesystemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *SharedFilesystemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SharedFilesystemModel
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
	if !plan.EncryptionEnabled.IsNull() && !plan.EncryptionEnabled.IsUnknown() {
		extraAttrs["encryption_enabled"] = plan.EncryptionEnabled.ValueBool()
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "shared_filesystem", plan.FilesystemName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.Protocol.IsUnknown() || plan.Protocol.IsNull() {
		plan.Protocol = types.StringValue("NFSv4")
	}
	if plan.PerformanceMode.IsUnknown() || plan.PerformanceMode.IsNull() {
		plan.PerformanceMode = types.StringValue("generalPurpose")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SharedFilesystemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SharedFilesystemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "shared_filesystem", state.FilesystemName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SharedFilesystemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SharedFilesystemModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "shared_filesystem", plan.FilesystemName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SharedFilesystemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SharedFilesystemModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "shared_filesystem", state.FilesystemName.ValueString(), state.ExtraConfig, resp)
}

func (r *SharedFilesystemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
