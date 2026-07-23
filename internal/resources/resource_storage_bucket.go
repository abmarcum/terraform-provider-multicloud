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

var _ resource.Resource = &StorageBucketResource{}
var _ resource.ResourceWithImportState = &StorageBucketResource{}

type StorageBucketResource struct {
	clientManager interface{}
}

type StorageBucketModel struct {
	ID                types.String `tfsdk:"id"`
	BucketName        types.String `tfsdk:"bucket_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	VersioningEnabled types.Bool   `tfsdk:"versioning_enabled"`
	EncryptionEnabled types.Bool   `tfsdk:"encryption_enabled"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewStorageBucketResource() resource.Resource {
	return &StorageBucketResource{}
}

func (r *StorageBucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_bucket"
}

func (r *StorageBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Storage Bucket resource supporting AWS S3, GCP Cloud Storage, and Azure Storage Container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "Unique identifier for the storage bucket resource.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the storage bucket across clouds.",
			},
			"provider_type": schema.StringAttribute{
				Required:    true,
				Description: "Target cloud provider ('aws', 'gcp', or 'azure').",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Cloud region for placement.",
			},
			"versioning_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable object versioning.",
			},
			"encryption_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Enable server-side encryption.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *StorageBucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *StorageBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageBucketModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "storage_bucket", plan.BucketName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.BucketName.IsUnknown() {
		if val, ok := res.Attributes["bucketname"].(string); ok && val != "" {
			plan.BucketName = types.StringValue(val)
		} else {
			plan.BucketName = types.StringValue("default-bucketname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
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

func (r *StorageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageBucketModel
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

	resName := state.BucketName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "storage_bucket", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageBucketModel
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

	resName := plan.BucketName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "storage_bucket", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "storage_bucket", state.BucketName.ValueString(), reg)
}

func (r *StorageBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
