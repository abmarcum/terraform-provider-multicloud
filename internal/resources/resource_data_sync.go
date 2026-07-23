package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DataSyncResource{}
var _ resource.ResourceWithImportState = &DataSyncResource{}

type DataSyncResource struct {
	clientManager interface{}
}

type DataSyncModel struct {
	ID                types.String `tfsdk:"id"`
	SyncName          types.String `tfsdk:"sync_name"`
	SourceProvider    types.String `tfsdk:"source_provider"`
	SourceBucket      types.String `tfsdk:"source_bucket"`
	DestProvider      types.String `tfsdk:"destination_provider"`
	DestBucket        types.String `tfsdk:"destination_bucket"`
	SyncSchedule      types.String `tfsdk:"sync_schedule"`
	SyncStatus        types.String `tfsdk:"sync_status"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"` 
	Region       types.String `tfsdk:"region"` 
}

func NewDataSyncResource() resource.Resource {
	return &DataSyncResource{}
}

func (r *DataSyncResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_sync"
}

func (r *DataSyncResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Data Synchronization resource replicating storage objects between AWS S3, GCP Cloud Storage, and Azure Blob Storage.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"sync_name": schema.StringAttribute{
				Required: true,
			},
			"source_provider": schema.StringAttribute{
				Required: true,
			},
			"source_bucket": schema.StringAttribute{
				Required: true,
			},
			"destination_provider": schema.StringAttribute{
				Required: true,
			},
			"destination_bucket": schema.StringAttribute{
				Required: true,
			},
			"sync_schedule": schema.StringAttribute{
				Optional: true,
			},
			"sync_status": schema.StringAttribute{
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

func (r *DataSyncResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *DataSyncResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DataSyncModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.SyncName.ValueString()
	src := strings.ToLower(plan.SourceProvider.ValueString())
	dst := strings.ToLower(plan.DestProvider.ValueString())

	plan.ID = types.StringValue(fmt.Sprintf("datasync/%s/%s-to-%s", name, src, dst))
	plan.SyncStatus = types.StringValue("REPLICATION_ACTIVE")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataSyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DataSyncModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DataSyncResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DataSyncModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DataSyncResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DataSyncModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataSyncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
