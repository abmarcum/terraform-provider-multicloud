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

var _ resource.Resource = &StorageInventoryReportResource{}
var _ resource.ResourceWithImportState = &StorageInventoryReportResource{}

type StorageInventoryReportResource struct {
	clientManager interface{}
}

type StorageInventoryReportModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	BucketName        types.String `tfsdk:"bucket_name"`
	DestinationBucket types.String `tfsdk:"destination_bucket"`
	Format            types.String `tfsdk:"format"`
	ScheduleFrequency types.String `tfsdk:"schedule_frequency"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewStorageInventoryReportResource() resource.Resource {
	return &StorageInventoryReportResource{}
}

func (r *StorageInventoryReportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_inventory_report"
}

func (r *StorageInventoryReportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Storage Inventory Report resource supporting AWS S3 Bucket Inventory, GCP Cloud Storage Inventory Reports, and Azure Storage Blob Inventory.",
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
			"bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "Target object storage bucket to audit and generate inventory reports for.",
			},
			"destination_bucket": schema.StringAttribute{
				Required:    true,
				Description: "Destination storage bucket where inventory report files (CSV/ORC/Parquet) are delivered.",
			},
			"format": schema.StringAttribute{
				Optional:    true,
				Description: "Inventory report file format: 'CSV' (default), 'PARQUET', or 'ORC'.",
			},
			"schedule_frequency": schema.StringAttribute{
				Optional:    true,
				Description: "Report generation frequency: 'DAILY' (default) or 'WEEKLY'.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *StorageInventoryReportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *StorageInventoryReportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageInventoryReportModel
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
	extraAttrs["bucket_name"] = plan.BucketName.ValueString()
	extraAttrs["destination_bucket"] = plan.DestinationBucket.ValueString()

	res, err := adapters.CreateCloudResource(ctx, providerType, "storage_inventory_report", plan.Name.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Name.IsUnknown() {
		plan.Name = types.StringValue("default-inventory-report")
	}
	if plan.ProviderType.IsUnknown() {
		plan.ProviderType = types.StringValue(providerType)
	}
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.Format.IsUnknown() || plan.Format.IsNull() {
		plan.Format = types.StringValue("CSV")
	}
	if plan.ScheduleFrequency.IsUnknown() || plan.ScheduleFrequency.IsNull() {
		plan.ScheduleFrequency = types.StringValue("DAILY")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageInventoryReportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageInventoryReportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "storage_inventory_report", state.Name.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageInventoryReportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageInventoryReportModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "storage_inventory_report", plan.Name.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageInventoryReportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageInventoryReportModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "storage_inventory_report", state.Name.ValueString(), state.ExtraConfig, resp)
}

func (r *StorageInventoryReportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
