package security

import (
	"testing"
)

func TestSecurityAuditorMatrix(t *testing.T) {
	// 1. Storage bucket unencrypted + public exposure (2 findings)
	findings := AuditResource("azure", "multicloud_storage_bucket", "unencrypted-public-container", true, false)
	if len(findings) != 2 {
		t.Fatalf("expected 2 audit findings, got %d", len(findings))
	}

	// 2. Storage bucket encrypted + private (0 findings)
	findings = AuditResource("aws", "multicloud_storage_bucket", "secure-private-bucket", false, true)
	if len(findings) != 0 {
		t.Fatalf("expected 0 audit findings for secure bucket, got %d", len(findings))
	}
}
