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

var _ resource.Resource = &IdentityFederationResource{}
var _ resource.ResourceWithImportState = &IdentityFederationResource{}

type IdentityFederationResource struct {
	clientManager interface{}
}

type IdentityFederationModel struct {
	ID              types.String `tfsdk:"id"`
	FederationName  types.String `tfsdk:"federation_name"`
	IssuerCloud     types.String `tfsdk:"issuer_cloud"`
	SubjectCloud    types.String `tfsdk:"subject_cloud"`
	Audience        types.String `tfsdk:"audience"`
	TrustStatus     types.String `tfsdk:"trust_status"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"` 
	Region       types.String `tfsdk:"region"` 
}

func NewIdentityFederationResource() resource.Resource {
	return &IdentityFederationResource{}
}

func (r *IdentityFederationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_federation"
}

func (r *IdentityFederationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Workload Identity Federation resource establishing keyless OIDC trust between AWS IAM, GCP Workload Identity, and Azure Entra ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"federation_name": schema.StringAttribute{
				Required: true,
			},
			"issuer_cloud": schema.StringAttribute{
				Required: true,
			},
			"subject_cloud": schema.StringAttribute{
				Required: true,
			},
			"audience": schema.StringAttribute{
				Optional: true,
			},
			"trust_status": schema.StringAttribute{
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

func (r *IdentityFederationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *IdentityFederationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IdentityFederationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	providerType := strings.ToLower(plan.IssuerCloud.ValueString())
	reg := ""
	if !plan.Region.IsNull() && !plan.Region.IsUnknown() {
		reg = plan.Region.ValueString()
	} else {
		plan.Region = types.StringNull()
	}
	res, err := adapters.CreateCloudResource(ctx, providerType, "identity_federation", plan.FederationName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.FederationName.IsUnknown() {
		if val, ok := res.Attributes["federationname"].(string); ok && val != "" {
			plan.FederationName = types.StringValue(val)
		} else {
			plan.FederationName = types.StringValue("default-federationname")
		}
	}
	if plan.IssuerCloud.IsUnknown() {
		if val, ok := res.Attributes["issuercloud"].(string); ok && val != "" {
			plan.IssuerCloud = types.StringValue(val)
		} else {
			plan.IssuerCloud = types.StringValue("default-issuercloud")
		}
	}
	if plan.SubjectCloud.IsUnknown() {
		if val, ok := res.Attributes["subjectcloud"].(string); ok && val != "" {
			plan.SubjectCloud = types.StringValue(val)
		} else {
			plan.SubjectCloud = types.StringValue("default-subjectcloud")
		}
	}
	if plan.Audience.IsUnknown() {
		if val, ok := res.Attributes["audience"].(string); ok && val != "" {
			plan.Audience = types.StringValue(val)
		} else {
			plan.Audience = types.StringValue("default-audience")
		}
	}
	if plan.TrustStatus.IsUnknown() {
		if val, ok := res.Attributes["truststatus"].(string); ok && val != "" {
			plan.TrustStatus = types.StringValue(val)
		} else {
			plan.TrustStatus = types.StringValue("default-truststatus")
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

func (r *IdentityFederationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IdentityFederationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.IssuerCloud, state.Region, "identity_federation", state.FederationName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IdentityFederationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IdentityFederationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.IssuerCloud, plan.Region, "identity_federation", plan.FederationName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IdentityFederationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IdentityFederationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.IssuerCloud, state.Region, "identity_federation", state.FederationName.ValueString(), state.ExtraConfig, resp)
}

func (r *IdentityFederationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
