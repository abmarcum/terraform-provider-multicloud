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

var _ resource.Resource = &DNSZoneLinkResource{}
var _ resource.ResourceWithImportState = &DNSZoneLinkResource{}

type DNSZoneLinkResource struct {
	clientManager interface{}
}

type DNSZoneLinkModel struct {
	ID                  types.String `tfsdk:"id"`
	LinkName            types.String `tfsdk:"link_name"`
	ZoneID              types.String `tfsdk:"zone_id"`
	VPCID               types.String `tfsdk:"vpc_id"`
	ProviderType        types.String `tfsdk:"provider_type"`
	RegistrationEnabled types.Bool   `tfsdk:"registration_enabled"`
	Region              types.String `tfsdk:"region"`
	ExtraConfig         types.Map    `tfsdk:"extra_config"`
}

func NewDNSZoneLinkResource() resource.Resource {
	return &DNSZoneLinkResource{}
}

func (r *DNSZoneLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_zone_link"
}

func (r *DNSZoneLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Private DNS Zone Virtual Network Attachment supporting AWS Route53 VPC Association, GCP Private DNS Visibility, and Azure Private DNS VNet Link.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"link_name": schema.StringAttribute{
				Required: true,
			},
			"zone_id": schema.StringAttribute{
				Required: true,
			},
			"vpc_id": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"registration_enabled": schema.BoolAttribute{
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

func (r *DNSZoneLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *DNSZoneLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSZoneLinkModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "dns_zone_link", plan.LinkName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
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

func (r *DNSZoneLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSZoneLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_zone_link", state.LinkName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSZoneLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSZoneLinkModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "dns_zone_link", plan.LinkName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSZoneLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSZoneLinkModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "dns_zone_link", state.LinkName.ValueString(), state.ExtraConfig, resp)
}

func (r *DNSZoneLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
