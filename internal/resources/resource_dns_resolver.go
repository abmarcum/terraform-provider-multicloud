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

var _ resource.Resource = &DNSResolverResource{}
var _ resource.ResourceWithImportState = &DNSResolverResource{}

type DNSResolverResource struct {
	clientManager interface{}
}

type DNSResolverModel struct {
	ID               types.String `tfsdk:"id"`
	ResolverName     types.String `tfsdk:"resolver_name"`
	ProviderType     types.String `tfsdk:"provider_type"`
	Direction        types.String `tfsdk:"direction"`
	VPCID            types.String `tfsdk:"vpc_id"`
	IPConfigurations types.List   `tfsdk:"ip_configurations"`
	Region           types.String `tfsdk:"region"`
	ExtraConfig      types.Map    `tfsdk:"extra_config"`
}

func NewDNSResolverResource() resource.Resource {
	return &DNSResolverResource{}
}

func (r *DNSResolverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_resolver"
}

func (r *DNSResolverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Hybrid DNS Resolver Endpoint supporting AWS Route53 Resolver Endpoint, GCP Cloud DNS Policy, and Azure Private DNS Resolver.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"resolver_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"direction": schema.StringAttribute{
				Required: true,
			},
			"vpc_id": schema.StringAttribute{
				Required: true,
			},
			"ip_configurations": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
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

func (r *DNSResolverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *DNSResolverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSResolverModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "dns_resolver", plan.ResolverName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
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

func (r *DNSResolverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSResolverModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_resolver", state.ResolverName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSResolverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSResolverModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "dns_resolver", plan.ResolverName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSResolverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSResolverModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_resolver", state.ResolverName.ValueString(), state.ExtraConfig, resp)
}

func (r *DNSResolverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
