package security

import (
	"fmt"
	"strings"
)

// SecurityAuditFinding represents a compliance audit result
type SecurityAuditFinding struct {
	ResourceName string
	ProviderType string
	Severity     string // "CRITICAL", "HIGH", "MEDIUM", "LOW"
	RuleID       string
	Message      string
}

// AuditResource evaluates a unified resource configuration against CIS Benchmarks & SOC 2 rules
func AuditResource(providerType string, resourceType string, resourceName string, isPublic bool, encryptionEnabled bool) []SecurityAuditFinding {
	var findings []SecurityAuditFinding
	p := strings.ToUpper(providerType)

	// CIS Rule 1.1: Disallow unencrypted storage
	if strings.Contains(resourceType, "storage") && !encryptionEnabled {
		findings = append(findings, SecurityAuditFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			Severity:     "HIGH",
			RuleID:       "CIS-STORAGE-ENCRYPTION-1.1",
			Message:      fmt.Sprintf("[%s CIS Audit] Storage resource '%s' must have server-side encryption enabled.", p, resourceName),
		})
	}

	// CIS Rule 2.4: Disallow public internet exposures
	if isPublic {
		findings = append(findings, SecurityAuditFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			Severity:     "CRITICAL",
			RuleID:       "CIS-NETWORK-NO-PUBLIC-EXPOSURE-2.4",
			Message:      fmt.Sprintf("[%s CIS Audit] Resource '%s' is publicly exposed to 0.0.0.0/0.", p, resourceName),
		})
	}

	return findings
}
