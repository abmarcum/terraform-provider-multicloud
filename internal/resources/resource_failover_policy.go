package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &FailoverPolicyResource{}
var _ resource.ResourceWithImportState = &FailoverPolicyResource{}

type FailoverPolicyResource struct {
	clientManager interface{}
}

type FailoverPolicyModel struct {
	ID             types.String `tfsdk:"id"`
	PolicyName     types.String `tfsdk:"policy_name"`
	PrimaryCloud   types.String `tfsdk:"primary_cloud"`
	FailoverCloud  types.String `tfsdk:"failover_cloud"`
	HealthCheckURL types.String `tfsdk:"health_check_url"`
	AutoFailover   types.Bool   `tfsdk:"auto_failover"`
	FailoverStatus types.String `tfsdk:"failover_status"`
	ExtraConfig    types.Map    `tfsdk:"extra_config"`
	Region       types.String `tfsdk:"region"` 
}

func NewFailoverPolicyResource() resource.Resource {
	return &FailoverPolicyResource{}
}

func (r *FailoverPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_failover_policy"
}

func (r *FailoverPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Disaster Recovery Failover Policy resource managing active-passive routing between AWS, GCP, and Azure.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"policy_name": schema.StringAttribute{
				Required: true,
			},
			"primary_cloud": schema.StringAttribute{
				Required:    true,
				Description: "Primary cloud target: 'aws', 'gcp', or 'azure'.",
			},
			"failover_cloud": schema.StringAttribute{
				Required:    true,
				Description: "Failover backup cloud target: 'aws', 'gcp', or 'azure'.",
			},
			"health_check_url": schema.StringAttribute{
				Optional: true,
			},
			"auto_failover": schema.BoolAttribute{
				Optional: true,
			},
			"failover_status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status: 'PRIMARY_HEALTHY' or 'FAILED_OVER'.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *FailoverPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *FailoverPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FailoverPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("failover/%s", plan.PolicyName.ValueString()))
	plan.FailoverStatus = types.StringValue("PRIMARY_HEALTHY")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FailoverPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FailoverPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FailoverPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FailoverPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *FailoverPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FailoverPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FailoverPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
