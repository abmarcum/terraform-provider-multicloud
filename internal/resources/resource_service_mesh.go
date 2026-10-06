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

var _ resource.Resource = &ServiceMeshResource{}
var _ resource.ResourceWithImportState = &ServiceMeshResource{}

type ServiceMeshResource struct {
	clientManager interface{}
}

type ServiceMeshModel struct {
	ID            types.String `tfsdk:"id"`
	MeshName      types.String `tfsdk:"mesh_name"`
	ProviderType  types.String `tfsdk:"provider_type"`
	Region        types.String `tfsdk:"region"`
	MTLSMode      types.String `tfsdk:"mtls_mode"`
	EgressFilter  types.String `tfsdk:"egress_filter"`
	ExtraConfig   types.Map    `tfsdk:"extra_config"`
}

func NewServiceMeshResource() resource.Resource {
	return &ServiceMeshResource{}
}

func (r *ServiceMeshResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_mesh"
}

func (r *ServiceMeshResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Service Mesh resource supporting AWS App Mesh, GCP Cloud Service Mesh, and Azure Kubernetes Fleet Service Mesh.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"mesh_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the service mesh control plane.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"mtls_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Mutual TLS enforcement policy ('STRICT' or 'PERMISSIVE').",
			},
			"egress_filter": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Outbound traffic filter mode ('ALLOW_ALL' or 'DROP_ALL').",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *ServiceMeshResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *ServiceMeshResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceMeshModel
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

	res, err := adapters.CreateCloudResource(ctx, providerType, "service_mesh", plan.MeshName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.MTLSMode.IsUnknown() || plan.MTLSMode.IsNull() {
		plan.MTLSMode = types.StringValue("STRICT")
	}
	if plan.EgressFilter.IsUnknown() || plan.EgressFilter.IsNull() {
		plan.EgressFilter = types.StringValue("DROP_ALL")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceMeshResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceMeshModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "service_mesh", state.MeshName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceMeshResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceMeshModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "service_mesh", plan.MeshName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceMeshResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ServiceMeshModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "service_mesh", state.MeshName.ValueString(), state.ExtraConfig, resp)
}

func (r *ServiceMeshResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
