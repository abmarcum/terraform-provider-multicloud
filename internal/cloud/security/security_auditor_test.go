package security

import (
	"testing"
)

func TestAuditResource(t *testing.T) {
	// 1. Test unencrypted storage finding
	findings := AuditResource("aws", "multicloud_storage_bucket", "my-unencrypted-bucket", false, false)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for unencrypted storage, got %d", len(findings))
	}
	if findings[0].Severity != "HIGH" {
		t.Errorf("expected Severity 'HIGH', got '%s'", findings[0].Severity)
	}

	// 2. Test public exposure finding
	findings = AuditResource("gcp", "multicloud_virtual_machine", "my-public-vm", true, true)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for public exposure, got %d", len(findings))
	}
	if findings[0].Severity != "CRITICAL" {
		t.Errorf("expected Severity 'CRITICAL', got '%s'", findings[0].Severity)
	}
}
