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

var _ resource.Resource = &SecretRotatorResource{}
var _ resource.ResourceWithImportState = &SecretRotatorResource{}

type SecretRotatorResource struct {
	clientManager interface{}
}

type SecretRotatorModel struct {
	ID                   types.String `tfsdk:"id"`
	RotationName         types.String `tfsdk:"rotation_name"`
	ProviderType         types.String `tfsdk:"provider_type"`
	SecretID             types.String `tfsdk:"secret_id"`
	RotationIntervalDays types.Int64  `tfsdk:"rotation_interval_days"`
	RotationStatus       types.String `tfsdk:"rotation_status"`
}

func NewSecretRotatorResource() resource.Resource {
	return &SecretRotatorResource{}
}

func (r *SecretRotatorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret_rotator"
}

func (r *SecretRotatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Secret Rotator resource managing automated secret rotation across AWS Secrets Manager, GCP Secret Manager, and Azure Key Vault.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"rotation_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"secret_id": schema.StringAttribute{
				Required: true,
			},
			"rotation_interval_days": schema.Int64Attribute{
				Optional: true,
			},
			"rotation_status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *SecretRotatorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *SecretRotatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SecretRotatorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	providerType := strings.ToLower(plan.ProviderType.ValueString())
	name := plan.RotationName.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf("%s/rotator/%s", providerType, name))
	plan.RotationStatus = types.StringValue("ROTATION_ENABLED")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecretRotatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SecretRotatorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SecretRotatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SecretRotatorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SecretRotatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *SecretRotatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
