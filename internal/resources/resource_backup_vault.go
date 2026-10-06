package resources

import (
	"context"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &BackupVaultResource{}
var _ resource.ResourceWithImportState = &BackupVaultResource{}

type BackupVaultResource struct {
	clientManager interface{}
}

type BackupVaultModel struct {
	ID                types.String `tfsdk:"id"`
	VaultName         types.String `tfsdk:"vault_name"`
	ProviderType      types.String `tfsdk:"provider_type"`
	Region            types.String `tfsdk:"region"`
	RetentionDays     types.Int64  `tfsdk:"retention_days"`
	ImmutableLock     types.Bool   `tfsdk:"immutable_lock"`
	EncryptionEnabled types.Bool   `tfsdk:"encryption_enabled"`
	ExtraConfig       types.Map    `tfsdk:"extra_config"`
}

func NewBackupVaultResource() resource.Resource {
	return &BackupVaultResource{}
}

func (r *BackupVaultResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backup_vault"
}

func (r *BackupVaultResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Disaster Recovery Backup Vault resource supporting AWS Backup Vaults, GCP Backup & DR Service, and Azure Data Protection Backup Vaults.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"vault_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the disaster recovery backup vault.",
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"retention_days": schema.Int64Attribute{
				Optional:    true,
				Description: "Backup recovery point retention window in days.",
			},
			"immutable_lock": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable WORM (Write-Once-Read-Many) immutable backup vault lock.",
			},
			"encryption_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Enable KMS encryption for backup vault recovery points.",
			},
			"extra_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Cloud-specific escape hatch key-value parameters passed to upstream cloud SDKs.",
			},
		},
	}
}

func (r *BackupVaultResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *BackupVaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan BackupVaultModel
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
	if !plan.EncryptionEnabled.IsNull() && !plan.EncryptionEnabled.IsUnknown() {
		extraAttrs["encryption_enabled"] = plan.EncryptionEnabled.ValueBool()
	}

	res, err := adapters.CreateCloudResource(ctx, providerType, "backup_vault", plan.VaultName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, extraAttrs))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.Region.IsUnknown() {
		plan.Region = types.StringValue("us-central1")
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupVaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state BackupVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "backup_vault", state.VaultName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BackupVaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan BackupVaultModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "backup_vault", plan.VaultName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *BackupVaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state BackupVaultModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "backup_vault", state.VaultName.ValueString(), state.ExtraConfig, resp)
}

func (r *BackupVaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
