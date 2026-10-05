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
	return EvaluateOPARegoPolicyWithAttrs(providerType, resourceType, resourceName, regoPolicyRules, nil)
}

// EvaluateOPARegoPolicyWithAttrs evaluates resource configuration and attributes against custom Rego policy rules
func EvaluateOPARegoPolicyWithAttrs(providerType string, resourceType string, resourceName string, regoPolicyRules string, attrs map[string]interface{}) OPARuleResult {
	p := strings.ToLower(providerType)
	r := strings.ToLower(resourceType)

	// Rego Rule 1: Mandatory Environment Tagging
	if strings.Contains(regoPolicyRules, "must_have_tags") {
		if attrs != nil {
			hasTags := false
			switch t := attrs["tags"].(type) {
			case map[string]string:
				hasTags = len(t) > 0
			case map[string]interface{}:
				hasTags = len(t) > 0
			}
			if !hasTags {
				return OPARuleResult{
					ResourceName: resourceName,
					ProviderType: providerType,
					PolicyRule:   "rego.mandatory_tagging",
					Passed:       false,
					Violation:    fmt.Sprintf("[OPA Rego Policy Violation] Resource '%s' on %s is missing mandatory tags.", resourceName, strings.ToUpper(p)),
				}
			}
		}
		return OPARuleResult{
			ResourceName: resourceName,
			ProviderType: providerType,
			PolicyRule:   "rego.mandatory_tagging",
			Passed:       true,
			Violation:    "",
		}
	}

	// Rego Rule 2: Multi-Cloud Region / Public IP Restrictions
	if strings.Contains(regoPolicyRules, "disallow_public_ip") {
		publicRequested := true
		if attrs != nil {
			if pub, ok := attrs["associate_public_ip"].(bool); ok {
				publicRequested = pub
			} else if pub, ok := attrs["is_public"].(bool); ok {
				publicRequested = pub
			}
		}
		if strings.Contains(r, "virtual_machine") && publicRequested {
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
