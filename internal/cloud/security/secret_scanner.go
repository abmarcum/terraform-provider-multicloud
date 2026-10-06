package security

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	awsSecretKeyRegex      = regexp.MustCompile(`(?i)(aws_secret_access_key|aws_secret_key)\s*=\s*["']?[A-Za-z0-9/+=]{40}["']?`)
	awsAccessKeyIDRegex    = regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)
	privateKeyRegex        = regexp.MustCompile(`-----BEGIN (?:(?:RSA|EC|DSA|OPENSSH|ENCRYPTED) )?PRIVATE KEY-----`)
	gcpServiceAccountRegex = regexp.MustCompile(`"type":\s*"service_account"`)
	jwtTokenRegex          = regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`)
	githubTokenRegex       = regexp.MustCompile(`\b(?:ghp_[A-Za-z0-9]{36}|github_pat_[A-Za-z0-9_]{22,})\b`)
)

// SecretScanFinding models a detected credential leak risk
type SecretScanFinding struct {
	ResourceName string
	ProviderType string
	SecretType   string
	Severity     string
	Message      string
}

// ScanForSecretLeaks inspects resource state attributes to block hardcoded secret leaks pre-apply
func ScanForSecretLeaks(providerType string, resourceName string, content string) []SecretScanFinding {
	var findings []SecretScanFinding
	p := strings.ToUpper(providerType)

	if awsSecretKeyRegex.MatchString(content) || awsAccessKeyIDRegex.MatchString(content) {
		/* #nosec G101 */
		findings = append(findings, SecretScanFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			SecretType:   "AWS_SECRET_ACCESS_KEY",
			Severity:     "CRITICAL",
			Message:      fmt.Sprintf("[%s SecretScanner] Hardcoded AWS Access/Secret Key pattern detected in resource '%s'.", p, resourceName),
		})
	}

	if privateKeyRegex.MatchString(content) {
		/* #nosec G101 */
		findings = append(findings, SecretScanFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			SecretType:   "PRIVATE_KEY_PEM",
			Severity:     "CRITICAL",
			Message:      fmt.Sprintf("[%s SecretScanner] Hardcoded Private Key PEM detected in resource '%s'.", p, resourceName),
		})
	}

	if gcpServiceAccountRegex.MatchString(content) {
		/* #nosec G101 */
		findings = append(findings, SecretScanFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			SecretType:   "GCP_SERVICE_ACCOUNT_JSON",
			Severity:     "HIGH",
			Message:      fmt.Sprintf("[%s SecretScanner] Hardcoded GCP Service Account JSON key detected in resource '%s'.", p, resourceName),
		})
	}

	if jwtTokenRegex.MatchString(content) {
		/* #nosec G101 */
		findings = append(findings, SecretScanFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			SecretType:   "JWT_TOKEN",
			Severity:     "HIGH",
			Message:      fmt.Sprintf("[%s SecretScanner] Hardcoded JWT Token pattern detected in resource '%s'.", p, resourceName),
		})
	}

	if githubTokenRegex.MatchString(content) {
		/* #nosec G101 */
		findings = append(findings, SecretScanFinding{
			ResourceName: resourceName,
			ProviderType: providerType,
			SecretType:   "GITHUB_PAT",
			Severity:     "HIGH",
			Message:      fmt.Sprintf("[%s SecretScanner] Hardcoded GitHub Personal Access Token detected in resource '%s'.", p, resourceName),
		})
	}

	return findings
}
