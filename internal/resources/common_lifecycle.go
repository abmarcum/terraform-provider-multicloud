package resources

import (
	"context"
	"errors"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/adapters"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// isReservedExtraConfigKey blocks untrusted extra_config entries from overriding provider credentials or mock_mode.
func isReservedExtraConfigKey(k string) bool {
	return adapters.IsSensitiveOrInternalKey(k)
}

// buildResourceExtraAttrs merges provider-level default attributes from ClientManager,
// resource-level extra_config map entries, and any resource-specific attributes.
func buildResourceExtraAttrs(clientManager interface{}, providerType string, extraConfig types.Map, custom map[string]interface{}) map[string]interface{} {
	attrs := make(map[string]interface{})

	if !extraConfig.IsNull() && !extraConfig.IsUnknown() {
		for k, v := range extraConfig.Elements() {
			if isReservedExtraConfigKey(k) {
				continue
			}
			if strVal, ok := v.(types.String); ok && !strVal.IsNull() && !strVal.IsUnknown() {
				attrs[k] = strVal.ValueString()
			}
		}
	}

	for k, v := range custom {
		attrs[k] = v
	}

	if clientManager != nil {
		attrs["_client_manager"] = clientManager
		if cm, ok := clientManager.(adapters.ClientConfigProvider); ok {
			for k, v := range cm.DefaultAttrs(providerType) {
				attrs[k] = v
			}
		}
	}

	return attrs
}

// resolveProviderAndRegion determines the normalized cloud provider and region,
// falling back to provider-level configuration defaults when region is omitted.
func resolveProviderAndRegion(providerAttr, regionAttr types.String, clientManager interface{}) (string, string) {
	pType := "gcp"
	if !providerAttr.IsNull() && !providerAttr.IsUnknown() && strings.TrimSpace(providerAttr.ValueString()) != "" {
		pType = strings.ToLower(strings.TrimSpace(providerAttr.ValueString()))
	}

	reg := ""
	if !regionAttr.IsNull() && !regionAttr.IsUnknown() && strings.TrimSpace(regionAttr.ValueString()) != "" {
		reg = strings.TrimSpace(regionAttr.ValueString())
	}

	if reg == "" && clientManager != nil {
		if cm, ok := clientManager.(adapters.ClientConfigProvider); ok {
			defaults := cm.DefaultAttrs(pType)
			if r, ok := defaults["provider_default_region"].(string); ok && r != "" {
				reg = r
			}
		}
	}

	if reg == "" {
		switch pType {
		case "aws":
			reg = "us-east-1"
		case "azure":
			reg = "eastus"
		default:
			reg = "us-central1"
		}
	}

	return pType, reg
}

// readResourceLifecycle performs a cloud resource read and handles ErrNotFound vs transient errors.
func readResourceLifecycle(
	ctx context.Context,
	clientManager interface{},
	providerAttr, regionAttr types.String,
	resType, resName string,
	extraConfig types.Map,
	resp *resource.ReadResponse,
) bool {
	pType, reg := resolveProviderAndRegion(providerAttr, regionAttr, clientManager)
	attrs := buildResourceExtraAttrs(clientManager, pType, extraConfig, nil)

	_, err := adapters.ReadCloudResourceWithAttrs(ctx, pType, resType, resName, reg, attrs)
	if err != nil {
		if errors.Is(err, adapters.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return false
		}
		resp.Diagnostics.AddError("Cloud Read Error", err.Error())
		return false
	}
	return true
}

// updateResourceLifecycle performs a cloud resource update with merged attributes.
func updateResourceLifecycle(
	ctx context.Context,
	clientManager interface{},
	providerAttr, regionAttr types.String,
	resType, resName string,
	extraConfig types.Map,
	custom map[string]interface{},
	resp *resource.UpdateResponse,
) bool {
	pType, reg := resolveProviderAndRegion(providerAttr, regionAttr, clientManager)
	attrs := buildResourceExtraAttrs(clientManager, pType, extraConfig, custom)

	_, err := adapters.UpdateCloudResource(ctx, pType, resType, resName, reg, attrs)
	if err != nil {
		resp.Diagnostics.AddError("Cloud Update Error", err.Error())
		return false
	}
	return true
}

// deleteResourceLifecycle performs a cloud resource deletion and surfaces errors unless already not found.
func deleteResourceLifecycle(
	ctx context.Context,
	clientManager interface{},
	providerAttr, regionAttr types.String,
	resType, resName string,
	extraConfig types.Map,
	resp *resource.DeleteResponse,
) {
	pType, reg := resolveProviderAndRegion(providerAttr, regionAttr, clientManager)
	attrs := buildResourceExtraAttrs(clientManager, pType, extraConfig, nil)

	err := adapters.DeleteCloudResourceWithAttrs(ctx, pType, resType, resName, reg, attrs)
	if err != nil && !errors.Is(err, adapters.ErrNotFound) {
		resp.Diagnostics.AddError("Cloud Delete Error", err.Error())
	}
}
