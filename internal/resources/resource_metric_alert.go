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

var _ resource.Resource = &MetricAlertResource{}
var _ resource.ResourceWithImportState = &MetricAlertResource{}

type MetricAlertResource struct {
	clientManager interface{}
}

type MetricAlertModel struct {
	ID           types.String  `tfsdk:"id"`
	AlertName    types.String  `tfsdk:"alert_name"`
	ProviderType types.String  `tfsdk:"provider_type"`
	MetricName   types.String  `tfsdk:"metric_name"`
	Threshold    types.Float64 `tfsdk:"threshold"`
	Comparison   types.String  `tfsdk:"comparison"`
	Region       types.String  `tfsdk:"region"`
	ExtraConfig  types.Map     `tfsdk:"extra_config"`
}

func NewMetricAlertResource() resource.Resource {
	return &MetricAlertResource{}
}

func (r *MetricAlertResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_metric_alert"
}

func (r *MetricAlertResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Multi-Cloud Metric Threshold Alarm Rule resource supporting AWS CloudWatch Alarm, GCP Monitoring Alert Policy, and Azure Metric Alert.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"alert_name": schema.StringAttribute{
				Required: true,
			},
			"provider_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"metric_name": schema.StringAttribute{
				Optional: true,
			},
			"threshold": schema.Float64Attribute{
				Optional: true,
			},
			"comparison": schema.StringAttribute{
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

func (r *MetricAlertResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientManager = req.ProviderData
}

func (r *MetricAlertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MetricAlertModel
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
	res, err := adapters.CreateCloudResource(ctx, providerType, "metric_alert", plan.AlertName.ValueString(), reg, buildResourceExtraAttrs(r.clientManager, providerType, plan.ExtraConfig, nil))
	if err != nil {
		resp.Diagnostics.AddError("Cloud Provision Error", err.Error())
		return
	}
	plan.ID = types.StringValue(res.ID)
	if plan.AlertName.IsUnknown() {
		if val, ok := res.Attributes["alertname"].(string); ok && val != "" {
			plan.AlertName = types.StringValue(val)
		} else {
			plan.AlertName = types.StringValue("default-alertname")
		}
	}
	if plan.ProviderType.IsUnknown() {
		if val, ok := res.Attributes["providertype"].(string); ok && val != "" {
			plan.ProviderType = types.StringValue(val)
		} else {
			plan.ProviderType = types.StringValue("default-providertype")
		}
	}
	if plan.MetricName.IsUnknown() {
		if val, ok := res.Attributes["metricname"].(string); ok && val != "" {
			plan.MetricName = types.StringValue(val)
		} else {
			plan.MetricName = types.StringValue("default-metricname")
		}
	}
	if plan.Comparison.IsUnknown() {
		if val, ok := res.Attributes["comparison"].(string); ok && val != "" {
			plan.Comparison = types.StringValue(val)
		} else {
			plan.Comparison = types.StringValue("default-comparison")
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

func (r *MetricAlertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MetricAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !readResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "metric_alert", state.AlertName.ValueString(), state.ExtraConfig, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MetricAlertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MetricAlertModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !updateResourceLifecycle(ctx, r.clientManager, plan.ProviderType, plan.Region, "metric_alert", plan.AlertName.ValueString(), plan.ExtraConfig, nil, resp) {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MetricAlertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MetricAlertModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	deleteResourceLifecycle(ctx, r.clientManager, state.ProviderType, state.Region, "metric_alert", state.AlertName.ValueString(), state.ExtraConfig, resp)
}

func (r *MetricAlertResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
