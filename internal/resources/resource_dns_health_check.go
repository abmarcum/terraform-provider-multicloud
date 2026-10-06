package resources

import (
	"context"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DNSHealthCheckResource{}
var _ resource.ResourceWithImportState = &DNSHealthCheckResource{}

type DNSHealthCheckResource struct {
	clientManager interface{}
}

type DNSHealthCheckModel struct {
	ID               types.String `tfsdk:"id"`
	CheckName        types.String `tfsdk:"check_name"`
	ProviderType     types.String `tfsdk:"provider_type"`
	Type             types.String `tfsdk:"type"`
	FQDN             types.String `tfsdk:"fqdn"`
	IPAddress        types.String `tfsdk:"ip_address"`
	Port             types.Int64  `tfsdk:"port"`
	ResourcePath     types.String `tfsdk:"resource_path"`
	FailureThreshold types.Int64  `tfsdk:"failure_threshold"`
	Region           types.String `tfsdk:"region"`
	ExtraConfig      types.Map    `tfsdk:"extra_config"`
}

func NewDNSHealthCheckResource() resource.Resource {
	return &DNSHealthCheckResource{}
}

func (r *DNSHealthCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_health_check"
}

func (r *DNSHealthCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Automated Endpoint Health Check resource supporting AWS Route53 Health Check, GCP Monitoring Uptime Probe, and Azure Traffic Manager Probe.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"check_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"fqdn": schema.StringAttribute{
				Optional: true,
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
			},
			"port": schema.Int64Attribute{
				Optional: true,
			},
			"resource_path": schema.StringAttribute{
				Optional: true,
			},
			"failure_threshold": schema.Int64Attribute{
				Optional: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (r *DNSHealthCheckResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *DNSHealthCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSHealthCheckModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "dns_health_check", plan.CheckName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSHealthCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSHealthCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_health_check", state.CheckName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSHealthCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSHealthCheckModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "dns_health_check", plan.CheckName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSHealthCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSHealthCheckModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_health_check", state.CheckName.ValueString(), state.ExtraConfig, resp)
}

func (r *DNSHealthCheckResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
