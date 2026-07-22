package security

import (
	"fmt"
	"strings"
)

// PolicyViolation represents an architectural policy evaluation failure
type PolicyViolation struct {
	RuleName     string
	ResourceName string
	Severity     string
	Message      string
}

// ValidatePolicy evaluates unified resource attributes against organizational security guidelines
func ValidatePolicy(resourceType string, resourceName string, attributes map[string]interface{}) []PolicyViolation {
	var violations []PolicyViolation

	// Rule 1: Storage buckets must not be public
	if strings.Contains(resourceType, "storage_bucket") {
		if public, ok := attributes["is_public"].(bool); ok && public {
			violations = append(violations, PolicyViolation{
				RuleName:     "NO_PUBLIC_STORAGE",
				ResourceName: resourceName,
				Severity:     "CRITICAL",
				Message:      fmt.Sprintf("Resource '%s' violates security policy: Storage buckets cannot be publicly accessible.", resourceName),
			})
		}
	}

	return violations
}
