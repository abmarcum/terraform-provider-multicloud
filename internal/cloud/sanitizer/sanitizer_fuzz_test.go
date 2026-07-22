package sanitizer

import (
	"testing"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
)

func FuzzSanitizeResourceName(f *testing.F) {
	// Seed initial test corpus
	f.Add("my-storage-bucket", "aws", "storage_bucket")
	f.Add("Azure_Account_123", "azure", "storage_bucket")
	f.Add("gcp_gcs_bucket", "gcp", "storage_bucket")

	f.Fuzz(func(t *testing.T, name string, provider string, resType string) {
		res := SanitizeResourceName(name, provider, resType)
		if len(res) == 0 && len(name) > 0 {
			// Ensure non-empty inputs return non-empty sanitized names
			t.Errorf("SanitizeResourceName(%q, %q, %q) returned empty string", name, provider, resType)
		}
	})
}

func FuzzScanForSecretLeaks(f *testing.F) {
	f.Add("aws", "bucket-1", "aws_secret_key = 'abcdef1234567890abcdef1234567890abcdef12'")
	f.Add("gcp", "vm-1", "-----BEGIN RSA PRIVATE KEY-----")

	f.Fuzz(func(t *testing.T, provider string, resName string, content string) {
		// Should not panic on arbitrary inputs
		_ = security.ScanForSecretLeaks(provider, resName, content)
	})
}
