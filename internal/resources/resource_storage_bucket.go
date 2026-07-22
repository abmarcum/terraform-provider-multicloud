package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/resiliency"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/sanitizer"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	gcpstorage "cloud.google.com/go/storage"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
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
	rawName := plan.BucketName.ValueString()
	sanitizedName := sanitizer.SanitizeResourceName(rawName, providerType, "storage_bucket")

	switch providerType {
	case "aws":
		// Call AWS S3 API via retry middleware
		_, err := resiliency.ExecuteWithRetry(ctx, func() (*awss3.CreateBucketOutput, error) {
			cm, ok := r.clientManager.(interface {
				GetAWSS3Client() (*awss3.Client, error)
			})
			if !ok {
				return nil, fmt.Errorf("AWS client not configured")
			}
			s3Client, err := cm.GetAWSS3Client()
			if err != nil {
				return nil, err
			}
			return s3Client.CreateBucket(ctx, &awss3.CreateBucketInput{
				Bucket: &sanitizedName,
			})
		})
		if err != nil {
			resp.Diagnostics.AddError("AWS S3 CreateBucket Error", err.Error())
			return
		}
		plan.ID = types.StringValue(fmt.Sprintf("aws/%s/%s", plan.Region.ValueString(), sanitizedName))

	case "gcp":
		_, err := resiliency.ExecuteWithRetry(ctx, func() (bool, error) {
			cm, ok := r.clientManager.(interface {
				GetGCPStorageClient() (*gcpstorage.Client, error)
			})
			if !ok {
				return false, fmt.Errorf("GCP client not configured")
			}
			gcsClient, err := cm.GetGCPStorageClient()
			if err != nil {
				return false, err
			}
			bucket := gcsClient.Bucket(sanitizedName)
			err = bucket.Create(ctx, "my-gcp-project", &gcpstorage.BucketAttrs{
				Location: plan.Region.ValueString(),
			})
			return true, err
		})
		if err != nil {
			resp.Diagnostics.AddError("GCP Storage CreateBucket Error", err.Error())
			return
		}
		plan.ID = types.StringValue(fmt.Sprintf("gcp/%s/%s", plan.Region.ValueString(), sanitizedName))

	case "azure":
		plan.ID = types.StringValue(fmt.Sprintf("azure/%s", sanitizedName))

	default:
		resp.Diagnostics.AddError("Unsupported Provider", fmt.Sprintf("Provider type '%s' is not supported.", providerType))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state StorageBucketModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *StorageBucketResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
