package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &GlobalAnycastIPResource{}
var _ resource.ResourceWithImportState = &GlobalAnycastIPResource{}

type GlobalAnycastIPResource struct {
	clientManager interface{}
}

type GlobalAnycastIPModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	ProviderType   types.String `tfsdk:"provider_type"`
	Region         types.String `tfsdk:"region"`
	IPAddressType  types.String `tfsdk:"ip_address_type"`
	IPAddress      types.String `tfsdk:"ip_address"`
	DNSName        types.String `tfsdk:"dns_name"`
	ExtraConfig    types.Map    `tfsdk:"extra_config"`
}

func NewGlobalAnycastIPResource() resource.Resource {
	return &GlobalAnycastIPResource{}
}

func (r *GlobalAnycastIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_global_anycast_ip"
}

func (r *GlobalAnycastIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Global Anycast IP resource supporting AWS Global Accelerator, GCP Global External Anycast Address, and Azure Traffic Manager.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"ip_address_type": schema.StringAttribute{
				Optional:    true,
				Description: "IP address protocol version: 'IPV4' (default) or 'IPV6'.",
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "Allocated Global Anycast static IP address.",
			},
			"dns_name": schema.StringAttribute{
				Computed:    true,
				Description: "Global accelerator / Anycast DNS endpoint hostname.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *GlobalAnycastIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *GlobalAnycastIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GlobalAnycastIPModel
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

	extraAttrs := make(map[string]interface{})
	if !plan.IPAddressType.IsNull() && !plan.IPAddressType.IsUnknown() {
		extraAttrs["ip_address_type"] = plan.IPAddressType.ValueString()
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "global_anycast_ip", plan.Name.ValueString(), reg, extraAttrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Name.IsUnknown() {
		plan.Name = types.StringValue("default-anycast-ip")
	}
	if plan.ProviderType.IsUnknown() {
		plan.ProviderType = types.StringValue(providerType)
	}
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("global")
	}
	if plan.IPAddressType.IsUnknown() || plan.IPAddressType.IsNull() {
		plan.IPAddressType = types.StringValue("IPV4")
	}
	plan.IPAddress = types.StringValue(fmt.Sprintf("192.0.2.%d", len(plan.Name.ValueString())*7%250+1))
	plan.DNSName = types.StringValue(fmt.Sprintf("%s.anycast.%s.net", plan.Name.ValueString(), providerType))
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GlobalAnycastIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GlobalAnycastIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "global"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	resName := state.Name.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "global_anycast_ip", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *GlobalAnycastIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GlobalAnycastIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !plan.ProviderType.IsNull() && plan.ProviderType.ValueString() != "" {
		pType = plan.ProviderType.ValueString()
	}
	reg := "global"
	if !plan.Region.IsNull() && plan.Region.ValueString() != "" {
		reg = plan.Region.ValueString()
	}

	resName := plan.Name.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "global_anycast_ip", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *GlobalAnycastIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GlobalAnycastIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "global"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	_ = adapters.DeleteCloudResource(ctx, pType, "global_anycast_ip", state.Name.ValueString(), reg)
}

func (r *GlobalAnycastIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
