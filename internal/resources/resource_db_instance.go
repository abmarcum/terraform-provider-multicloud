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

var _ resource.Resource = &DBInstanceResource{}
var _ resource.ResourceWithImportState = &DBInstanceResource{}

type DBInstanceResource struct {
	clientManager interface{}
}

type DBInstanceModel struct {
	ID                  types.String `tfsdk:"id"`
	InstanceName        types.String `tfsdk:"instance_name"`
	ProviderType        types.String `tfsdk:"provider_type"`
	Engine              types.String `tfsdk:"engine"`
	EngineVersion       types.String `tfsdk:"engine_version"`
	SizeTier            types.String `tfsdk:"size_tier"`
	StorageGB           types.Int64  `tfsdk:"storage_gb"`
	MultiAZ             types.Bool   `tfsdk:"multi_az"`
	BackupRetentionDays types.Int64  `tfsdk:"backup_retention_days"`
	Username            types.String `tfsdk:"username"`
	Password            types.String `tfsdk:"password"`
	Region              types.String `tfsdk:"region"`
	ExtraConfig         types.Map    `tfsdk:"extra_config"`
}

func NewDBInstanceResource() resource.Resource {
	return &DBInstanceResource{}
}

func (r *DBInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_db_instance"
}

func (r *DBInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Managed Relational DB instance supporting AWS RDS, GCP Cloud SQL, and Azure Database.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"instance_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"engine": schema.StringAttribute{
				Optional:    true,
				Description: "Database engine: 'postgres' or 'mysql'.",
			},
			"engine_version": schema.StringAttribute{
				Optional: true,
			},
			"size_tier": schema.StringAttribute{
				Optional: true,
			},
			"storage_gb": schema.Int64Attribute{
				Optional: true,
			},
			"multi_az": schema.BoolAttribute{
				Optional: true,
			},
			"backup_retention_days": schema.Int64Attribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
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

func (r *DBInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *DBInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DBInstanceModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "db_instance", plan.InstanceName.ValueString(), reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.InstanceName.IsUnknown() {
		if val, ok := res.Attributes["instancename"].(string); ok && val != "" {
			plan.InstanceName = types.StringValue(val)
		} else {
			plan.InstanceName = types.StringValue("default-instancename")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.Engine.IsUnknown() {
		if val, ok := res.Attributes["engine"].(string); ok && val != "" {
			plan.Engine = types.StringValue(val)
		} else {
			plan.Engine = types.StringValue("default-engine")
		}
	}
	if plan.EngineVersion.IsUnknown() {
		if val, ok := res.Attributes["engineversion"].(string); ok && val != "" {
			plan.EngineVersion = types.StringValue(val)
		} else {
			plan.EngineVersion = types.StringValue("default-engineversion")
		}
	}
	if plan.SizeTier.IsUnknown() {
		if val, ok := res.Attributes["sizetier"].(string); ok && val != "" {
			plan.SizeTier = types.StringValue(val)
		} else {
			plan.SizeTier = types.StringValue("default-sizetier")
		}
	}
	if plan.Username.IsUnknown() {
		if val, ok := res.Attributes["username"].(string); ok && val != "" {
			plan.Username = types.StringValue(val)
		} else {
			plan.Username = types.StringValue("default-username")
		}
	}
	if plan.Password.IsUnknown() {
		if val, ok := res.Attributes["password"].(string); ok && val != "" {
			plan.Password = types.StringValue(val)
		} else {
			plan.Password = types.StringValue("default-password")
		}
	}
	if plan.Region.IsUnknown() {
		if val, ok := res.Attributes["region"].(string); ok && val != "" {
			plan.Region = types.StringValue(val)
		} else {
			plan.Region = types.StringValue("default-region")
		}
	}
	if plan.ExtraConfig.IsUnknown() {
		plan.ExtraConfig = types.MapNull(types.StringType)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DBInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DBInstanceModel
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

	resName := state.InstanceName.ValueString()
	_, err := adapters.ReadCloudResource(ctx, pType, "db_instance", resName, reg)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DBInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DBInstanceModel
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

	resName := plan.InstanceName.ValueString()
	_, err := adapters.UpdateCloudResource(ctx, pType, "db_instance", resName, reg, nil)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DBInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DBInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	pType := strings.ToLower(state.ProviderType.ValueString())
	reg := ""
	if !state.Region.IsNull() && !state.Region.IsUnknown() {
		reg = state.Region.ValueString()
	}
	_ = adapters.DeleteCloudResource(ctx, pType, "db_instance", state.InstanceName.ValueString(), reg)
}

func (r *DBInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
