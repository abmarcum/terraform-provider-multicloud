package security

import (
	"testing"
)

func TestScanForSecretLeaks(t *testing.T) {
	// 1. Test AWS Secret Key Leak
	findings := ScanForSecretLeaks("aws", "my-secret-resource", "aws_secret_key = \"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY\"")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for AWS secret key leak, got %d", len(findings))
	}
	if findings[0].Severity != "CRITICAL" {
		t.Errorf("expected Severity 'CRITICAL', got '%s'", findings[0].Severity)
	}

	// 2. Test Private Key PEM Leak
	findings = ScanForSecretLeaks("gcp", "my-pem-resource", "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQ...")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for RSA private key leak, got %d", len(findings))
	}

	// 3. Test Clean Input
	findings = ScanForSecretLeaks("azure", "clean-resource", "environment = 'production'")
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for clean content, got %d", len(findings))
	}
}
