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
	providerType := strings.ToLower(plan.SourceProvider.ValueString())
	reg := ""
	if !plan.Region.IsNull() && !plan.Region.IsUnknown() {
		reg = plan.Region.ValueString()
	} else {
		plan.Region = types.StringNull()
	}
	res, err := adapters.CreateCloudResource(ctx, providerType, "data_sync", plan.SyncName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.SyncName.IsUnknown() {
		if val, ok := res.Attributes["syncname"].(string); ok && val != "" {
			plan.SyncName = types.StringValue(val)
		} else {
			plan.SyncName = types.StringValue("default-syncname")
		}
	}
	if plan.SourceProvider.IsUnknown() {
		if val, ok := res.Attributes["sourceprovider"].(string); ok && val != "" {
			plan.SourceProvider = types.StringValue(val)
		} else {
			plan.SourceProvider = types.StringValue("default-sourceprovider")
		}
	}
	if plan.SourceBucket.IsUnknown() {
		if val, ok := res.Attributes["sourcebucket"].(string); ok && val != "" {
			plan.SourceBucket = types.StringValue(val)
		} else {
			plan.SourceBucket = types.StringValue("default-sourcebucket")
		}
	}
	if plan.DestProvider.IsUnknown() {
		if val, ok := res.Attributes["destprovider"].(string); ok && val != "" {
			plan.DestProvider = types.StringValue(val)
		} else {
			plan.DestProvider = types.StringValue("default-destprovider")
		}
	}
	if plan.DestBucket.IsUnknown() {
		if val, ok := res.Attributes["destbucket"].(string); ok && val != "" {
			plan.DestBucket = types.StringValue(val)
		} else {
			plan.DestBucket = types.StringValue("default-destbucket")
		}
	}
	if plan.SyncSchedule.IsUnknown() {
		if val, ok := res.Attributes["syncschedule"].(string); ok && val != "" {
			plan.SyncSchedule = types.StringValue(val)
		} else {
			plan.SyncSchedule = types.StringValue("default-syncschedule")
		}
	}
	if plan.SyncStatus.IsUnknown() {
		if val, ok := res.Attributes["syncstatus"].(string); ok && val != "" {
			plan.SyncStatus = types.StringValue(val)
		} else {
			plan.SyncStatus = types.StringValue("default-syncstatus")
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

func (r *DataSyncResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DataSyncModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.SourceProvider.IsNull() && state.SourceProvider.ValueString() != "" {
		pType = state.SourceProvider.ValueString()
	}
	reg := "us-central1"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	resName := state.SyncName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "data_sync", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
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

	pType := "gcp"
	if !plan.SourceProvider.IsNull() && plan.SourceProvider.ValueString() != "" {
		pType = plan.SourceProvider.ValueString()
	}
	reg := "us-central1"
	if !plan.Region.IsNull() && plan.Region.ValueString() != "" {
		reg = plan.Region.ValueString()
	}

	resName := plan.SyncName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "data_sync", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
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
	pType := strings.ToLower(state.SourceProvider.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "data_sync", state.SyncName.ValueString(), reg)
}

func (r *DataSyncResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
