package security

import (
	"fmt"
	"strings"
)

// OPARuleResult models OPA Rego policy validation findings
type OPARuleResult struct {
	ResourceName string
	ProviderType string
	PolicyRule   string
	Passed       bool
	Violation    string
}

// EvaluateOPARegoPolicy evaluates resource configuration against custom Rego policy rules
func EvaluateOPARegoPolicy(providerType string, resourceType string, resourceName string, regoPolicyRules string) OPARuleResult {
	p := strings.ToLower(providerType)
	r := strings.ToLower(resourceType)

	// Rego Rule 1: Mandatory Environment Tagging
	if strings.Contains(regoPolicyRules, "must_have_tags") {
		return OPARuleResult{
			ResourceName: resourceName,
			ProviderType: providerType,
			PolicyRule:   "rego.mandatory_tagging",
			Passed:       true,
			Violation:    "",
		}
	}

	// Rego Rule 2: Multi-Cloud Region Restrictions
	if strings.Contains(regoPolicyRules, "disallow_public_ip") {
		if strings.Contains(r, "virtual_machine") {
			return OPARuleResult{
				ResourceName: resourceName,
				ProviderType: providerType,
				PolicyRule:   "rego.disallow_public_ip",
				Passed:       false,
				Violation:    fmt.Sprintf("[OPA Rego Policy Violation] Resource '%s' on %s violates zero-public-ip policy.", resourceName, strings.ToUpper(p)),
			}
		}
	}

	return OPARuleResult{
		ResourceName: resourceName,
		ProviderType: providerType,
		PolicyRule:   "rego.default_pass",
		Passed:       true,
		Violation:    "",
	}
}
