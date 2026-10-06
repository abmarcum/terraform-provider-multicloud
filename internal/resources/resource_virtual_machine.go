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
	ID                types.String `tfsdk:"id"`
	VMName            types.String `tfsdk:"vm_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	SizeTier          types.String `tfsdk:"size_tier"`
	InstanceType      types.String `tfsdk:"instance_type"`
	ImageID           types.String `tfsdk:"image_id"`
	SubnetID          types.String `tfsdk:"subnet_id"`
	SSHPublicKey      types.String `tfsdk:"ssh_public_key"`
	OSDiskSizeGB      types.Int64  `tfsdk:"os_disk_size_gb"`
	OSDiskType        types.String `tfsdk:"os_disk_type"`
	UserData          types.String `tfsdk:"user_data"`
	AssociatePublicIP types.Bool   `tfsdk:"associate_public_ip"`
	Tags              types.Map    `tfsdk:"tags"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
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
				Description: "Standard instance tier: 'small', 'medium', 'large' (Defaults to Intel Xeon hardware).",
			},
			"instance_type": schema.StringAttribute{
				Optional:    true,
				Description: "Explicit cloud instance type/SKU specification (e.g., AWS 'm6i.xlarge', GCP 'n2-standard-4', Azure 'Standard_D4s_v5').",
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
			"os_disk_size_gb": schema.Int64Attribute{
				Optional: true,
			},
			"os_disk_type": schema.StringAttribute{
				Optional: true,
			},
			"user_data": schema.StringAttribute{
				Optional: true,
			},
			"associate_public_ip": schema.BoolAttribute{
				Optional: true,
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

	extraAttrs := make(map[string]interface{})
	if !plan.SizeTier.IsNull() && !plan.SizeTier.IsUnknown() {
		extraAttrs["size_tier"] = plan.SizeTier.ValueString()
	}
	if !plan.InstanceType.IsNull() && !plan.InstanceType.IsUnknown() && plan.InstanceType.ValueString() != "" {
		instType := plan.InstanceType.ValueString()
		extraAttrs["instance_type"] = instType
		extraAttrs["aws_instance_type"] = instType
		extraAttrs["gcp_machine_type"] = instType
		extraAttrs["azure_vm_sku"] = instType
	}
	if !plan.UserData.IsNull() && !plan.UserData.IsUnknown() && plan.UserData.ValueString() != "" {
		extraAttrs["user_data"] = plan.UserData.ValueString()
	}
	if !plan.AssociatePublicIP.IsNull() && !plan.AssociatePublicIP.IsUnknown() {
		extraAttrs["associate_public_ip"] = plan.AssociatePublicIP.ValueBool()
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "virtual_machine", plan.VMName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
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
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "virtual_machine", state.VMName.ValueString(), state.ExtraConfig, resp) {
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
	extraAttrs := make(map[string]interface{})
	if !plan.UserData.IsNull() && !plan.UserData.IsUnknown() && plan.UserData.ValueString() != "" {
		extraAttrs["user_data"] = plan.UserData.ValueString()
	}
	if !plan.AssociatePublicIP.IsNull() && !plan.AssociatePublicIP.IsUnknown() {
		extraAttrs["associate_public_ip"] = plan.AssociatePublicIP.ValueBool()
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "virtual_machine", plan.VMName.ValueString(), plan.ExtraConfig, extraAttrs, resp) {
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
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "virtual_machine", state.VMName.ValueString(), state.ExtraConfig, resp)
}

func (r *VirtualMachineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
