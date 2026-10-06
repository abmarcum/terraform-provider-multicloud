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

var _ resource.Resource = &VectorIndexResource{}
var _ resource.ResourceWithImportState = &VectorIndexResource{}

type VectorIndexResource struct {
	clientManager interface{}
}

type VectorIndexModel struct {
	ID             types.String `tfsdk:"id"`
	IndexName      types.String `tfsdk:"index_name"`
	ProviderType   types.String `tfsdk:"provider_type"`
	Region         types.String `tfsdk:"region"`
	Dimensions     types.Int64  `tfsdk:"dimensions"`
	DistanceMetric types.String `tfsdk:"distance_metric"`
	ExtraConfig    types.Map    `tfsdk:"extra_config"`
}

func NewVectorIndexResource() resource.Resource {
	return &VectorIndexResource{}
}

func (r *VectorIndexResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vector_index"
}

func (r *VectorIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud AI Vector Database Index resource supporting AWS OpenSearch Serverless Vector Engine, GCP Vertex AI Vector Search, and Azure AI Search Vector Index.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"index_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the vector embedding search index.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"dimensions": schema.Int64Attribute{
				Required:    true,
				Description: "Vector embedding dimensionality (e.g., 768, 1536, 3072).",
			},
			"distance_metric": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Similarity distance metric ('COSINE', 'DOT_PRODUCT', 'EUCLIDEAN').",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *VectorIndexResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *VectorIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VectorIndexModel
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
		"dimensions": plan.Dimensions.ValueInt64(),
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "vector_index", plan.IndexName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.DistanceMetric.IsUnknown() || plan.DistanceMetric.IsNull() {
		plan.DistanceMetric = types.StringValue("COSINE")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VectorIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VectorIndexModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "vector_index", state.IndexName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VectorIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VectorIndexModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "vector_index", plan.IndexName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VectorIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VectorIndexModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "vector_index", state.IndexName.ValueString(), state.ExtraConfig, resp)
}

func (r *VectorIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
