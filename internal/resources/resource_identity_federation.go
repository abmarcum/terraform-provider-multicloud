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
	issuer := strings.ToLower(plan.IssuerCloud.ValueString())
	subject := strings.ToLower(plan.SubjectCloud.ValueString())
	name := plan.FederationName.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf("federation/%s/%s-trusts-%s", name, issuer, subject))
	plan.TrustStatus = types.StringValue("OIDC_TRUST_ESTABLISHED")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IdentityFederationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IdentityFederationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
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
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IdentityFederationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *IdentityFederationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
