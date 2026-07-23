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

var _ resource.Resource = &VirtualMachineResource{}
var _ resource.ResourceWithImportState = &VirtualMachineResource{}

type VirtualMachineResource struct {
	clientManager interface{}
}

type VirtualMachineModel struct {
	ID           types.String `tfsdk:"id"`
	VMName       types.String `tfsdk:"vm_name"`
	ProviderType types.String `tfsdk:"provider_type"`
	Region       types.String `tfsdk:"region"`
	SizeTier     types.String `tfsdk:"size_tier"`
	ImageID      types.String `tfsdk:"image_id"`
	SubnetID     types.String `tfsdk:"subnet_id"`
	SSHPublicKey types.String `tfsdk:"ssh_public_key"`
	Tags         types.Map    `tfsdk:"tags"`
	ExtraConfig  types.Map    `tfsdk:"extra_config"`
}

func NewVirtualMachineResource() resource.Resource {
	return &VirtualMachineResource{}
}

func (r *VirtualMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

func (r *VirtualMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Virtual Machine resource supporting AWS EC2, GCP Compute Engine, and Azure VMs under a single schema.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vm_name": schema.StringAttribute{
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
			"size_tier": schema.StringAttribute{
				Optional:    true,
				Description: "Standard instance tier: 'small', 'medium', 'large'.",
			},
			"image_id": schema.StringAttribute{
				Optional: true,
			},
			"subnet_id": schema.StringAttribute{
				Optional: true,
			},
			"ssh_public_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"tags": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *VirtualMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *VirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VirtualMachineModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "virtual_machine", plan.VMName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.VMName.IsUnknown() {
		if val, ok := res.Attributes["vmname"].(string); ok && val != "" {
			plan.VMName = types.StringValue(val)
		} else {
			plan.VMName = types.StringValue("default-vmname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.Region.IsUnknown() {
		if val, ok := res.Attributes["region"].(string); ok && val != "" {
			plan.Region = types.StringValue(val)
		} else {
			plan.Region = types.StringValue("default-region")
		}
	}
	if plan.SizeTier.IsUnknown() {
		if val, ok := res.Attributes["sizetier"].(string); ok && val != "" {
			plan.SizeTier = types.StringValue(val)
		} else {
			plan.SizeTier = types.StringValue("default-sizetier")
		}
	}
	if plan.ImageID.IsUnknown() {
		if val, ok := res.Attributes["imageid"].(string); ok && val != "" {
			plan.ImageID = types.StringValue(val)
		} else {
			plan.ImageID = types.StringValue("default-imageid")
		}
	}
	if plan.SubnetID.IsUnknown() {
		if val, ok := res.Attributes["subnetid"].(string); ok && val != "" {
			plan.SubnetID = types.StringValue(val)
		} else {
			plan.SubnetID = types.StringValue("default-subnetid")
		}
	}
	if plan.SSHPublicKey.IsUnknown() {
		if val, ok := res.Attributes["sshpublickey"].(string); ok && val != "" {
			plan.SSHPublicKey = types.StringValue(val)
		} else {
			plan.SSHPublicKey = types.StringValue("default-sshpublickey")
		}
	}
	if plan.Tags.IsUnknown() {
		plan.Tags = types.MapNull(types.StringType)
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !state.ProviderType.IsNull() && state.ProviderType.ValueString() != "" {
		pType = state.ProviderType.ValueString()
	}
	reg := "us-central1"
	if !state.Region.IsNull() && state.Region.ValueString() != "" {
		reg = state.Region.ValueString()
	}

	resName := state.VMName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "virtual_machine", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VirtualMachineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pType := "gcp"
	if !plan.ProviderType.IsNull() && plan.ProviderType.ValueString() != "" {
		pType = plan.ProviderType.ValueString()
	}
	reg := "us-central1"
	if !plan.Region.IsNull() && plan.Region.ValueString() != "" {
		reg = plan.Region.ValueString()
	}

	resName := plan.VMName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "virtual_machine", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "virtual_machine", state.VMName.ValueString(), reg)
}

func (r *VirtualMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
