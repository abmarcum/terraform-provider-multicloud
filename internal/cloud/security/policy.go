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
	if attributes == nil {
		return violations
	}

	// Rule 1: Storage buckets must not be public and must not explicitly disable encryption
	if strings.Contains(resourceType, "storage_bucket") {
		if public, ok := attributes["is_public"].(bool); ok && public {
			violations = append(violations, PolicyViolation{
				RuleName:     "NO_PUBLIC_STORAGE",
				ResourceName: resourceName,
				Severity:     "CRITICAL",
				Message:      fmt.Sprintf("Resource '%s' violates security policy: Storage buckets cannot be publicly accessible.", resourceName),
			})
		}
		if enc, ok := attributes["encryption_enabled"].(bool); ok && !enc {
			violations = append(violations, PolicyViolation{
				RuleName:     "STORAGE_ENCRYPTION_REQUIRED",
				ResourceName: resourceName,
				Severity:     "HIGH",
				Message:      fmt.Sprintf("Resource '%s' violates security policy: Storage buckets must enable server-side encryption.", resourceName),
			})
		}
	}

	// Rule 2: Relational databases cannot be publicly accessible
	if strings.Contains(resourceType, "db_instance") {
		if public, ok := attributes["publicly_accessible"].(bool); ok && public {
			violations = append(violations, PolicyViolation{
				RuleName:     "DB_PUBLIC_ACCESS_DISALLOWED",
				ResourceName: resourceName,
				Severity:     "CRITICAL",
				Message:      fmt.Sprintf("Resource '%s' violates security policy: Database instances cannot be publicly accessible.", resourceName),
			})
		}
	}

	// Rule 3: KMS keys must not explicitly disable key rotation when enforced
	if strings.Contains(resourceType, "kms_key") {
		if rot, ok := attributes["require_key_rotation"].(bool); ok && rot {
			if enabled, ok := attributes["enable_key_rotation"].(bool); ok && !enabled {
				violations = append(violations, PolicyViolation{
					RuleName:     "KMS_ROTATION_REQUIRED",
					ResourceName: resourceName,
					Severity:     "HIGH",
					Message:      fmt.Sprintf("Resource '%s' violates security policy: KMS keys must have automatic key rotation enabled.", resourceName),
				})
			}
		}
	}

	return violations
}
