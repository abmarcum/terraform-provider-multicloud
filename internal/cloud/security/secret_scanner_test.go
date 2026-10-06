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

	// 2. Test Private Key PEM Leak (RSA and PKCS#8)
	findings = ScanForSecretLeaks("gcp", "my-pem-resource", "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQ...")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for RSA private key leak, got %d", len(findings))
	}
	findings = ScanForSecretLeaks("gcp", "my-pkcs8-resource", "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhki...")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for PKCS#8 private key leak, got %d", len(findings))
	}

	// 3. Test AWS Access Key ID and GitHub PAT
	findings = ScanForSecretLeaks("aws", "akia-res", "export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for AWS Access Key ID leak, got %d", len(findings))
	}
	findings = ScanForSecretLeaks("aws", "gh-res", "token = ghp_1234567890abcdefghijklmnopqrstuvwxyz")
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding for GitHub PAT leak, got %d", len(findings))
	}

	// 4. Test Clean Input
	findings = ScanForSecretLeaks("azure", "clean-resource", "environment = 'production'")
	if len(findings) != 0 {
		t.Errorf("expected 0 findings for clean content, got %d", len(findings))
	}
}
