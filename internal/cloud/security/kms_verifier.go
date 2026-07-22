package security

import (
	"fmt"
	"strings"
)

// KMSComplianceResult models customer-managed key verification findings
type KMSComplianceResult struct {
	ResourceName string
	ProviderType string
	IsCMK        bool
	Rotation     bool
	Violation    string
}

// VerifyKMSCompliance validates customer-managed key usage and 90-day rotation
func VerifyKMSCompliance(providerType string, resourceName string, kmsKeyARN string, rotationEnabled bool) KMSComplianceResult {
	p := strings.ToUpper(providerType)

	if kmsKeyARN == "" {
		return KMSComplianceResult{
			ResourceName: resourceName,
			ProviderType: providerType,
			IsCMK:        false,
			Rotation:     false,
			Violation:    fmt.Sprintf("[%s KMS Verifier] Resource '%s' is using cloud default keys instead of a Customer-Managed KMS Key (CMK).", p, resourceName),
		}
	}

	if !rotationEnabled {
		return KMSComplianceResult{
			ResourceName: resourceName,
			ProviderType: providerType,
			IsCMK:        true,
			Rotation:     false,
			Violation:    fmt.Sprintf("[%s KMS Verifier] KMS Key '%s' on resource '%s' lacks 90-day automated key rotation.", p, kmsKeyARN, resourceName),
		}
	}

	return KMSComplianceResult{
		ResourceName: resourceName,
		ProviderType: providerType,
		IsCMK:        true,
		Rotation:     true,
		Violation:    "",
	}
}
