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

var _ resource.Resource = &StorageTransferJobResource{}
var _ resource.ResourceWithImportState = &StorageTransferJobResource{}

type StorageTransferJobResource struct {
	clientManager interface{}
}

type StorageTransferJobModel struct {
	ID                types.String `tfsdk:"id"`
	JobName           types.String `tfsdk:"job_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	SourceBucket      types.String `tfsdk:"source_bucket"`
	DestinationBucket types.String `tfsdk:"destination_bucket"`
	ScheduleStartTime types.String `tfsdk:"schedule_start_time"`
	OverwriteObjects  types.Bool   `tfsdk:"overwrite_objects"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewStorageTransferJobResource() resource.Resource {
	return &StorageTransferJobResource{}
}

func (r *StorageTransferJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_transfer_job"
}

func (r *StorageTransferJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Storage Transfer Job resource supporting GCP Storage Transfer Service Jobs, AWS DataSync Tasks, and Azure Storage Sync Tasks.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"job_name": schema.StringAttribute{
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
			"source_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Source object storage bucket or URL.",
			},
			"destination_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Destination object storage bucket or URL.",
			},
			"schedule_start_time": schema.StringAttribute{
				Optional:    true,
				Description: "ISO-8601 schedule start time or cron expression for recurring data transfer.",
			},
			"overwrite_objects": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to overwrite existing destination objects (default true).",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *StorageTransferJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *StorageTransferJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageTransferJobModel
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
	extraAttrs["source_bucket"] = plan.SourceBucket.ValueString()
	extraAttrs["destination_bucket"] = plan.DestinationBucket.ValueString()

	res, err := adapters.CreateCloudResource(ctx, providerType, "storage_transfer_job", plan.JobName.ValueString(), reg, extraAttrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.JobName.IsUnknown() {
		plan.JobName = types.StringValue("default-transfer-job")
	}
	if plan.ProviderType.IsUnknown() {
		plan.ProviderType = types.StringValue(providerType)
	}
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.OverwriteObjects.IsUnknown() || plan.OverwriteObjects.IsNull() {
		plan.OverwriteObjects = types.BoolValue(true)
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageTransferJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageTransferJobModel
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

	resName := state.JobName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "storage_transfer_job", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageTransferJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageTransferJobModel
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

	resName := plan.JobName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "storage_transfer_job", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageTransferJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageTransferJobModel
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

	_ = adapters.DeleteCloudResource(ctx, pType, "storage_transfer_job", state.JobName.ValueString(), reg)
}

func (r *StorageTransferJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
