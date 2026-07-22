package security

import (
	"testing"
)

func TestVerifyKMSCompliance(t *testing.T) {
	// 1. Test missing CMK key
	res := VerifyKMSCompliance("aws", "my-bucket", "", false)
	if res.IsCMK {
		t.Errorf("expected IsCMK to be false for empty key ARN")
	}

	// 2. Test missing key rotation
	res = VerifyKMSCompliance("gcp", "my-gcs-bucket", "projects/gcp/locations/global/keyRings/ring/cryptoKeys/key1", false)
	if res.Rotation {
		t.Errorf("expected Rotation to be false when disabled")
	}

	// 3. Test compliant KMS key
	res = VerifyKMSCompliance("azure", "my-vault-key", "https://keyvault.vault.azure.net/keys/key1", true)
	if res.Violation != "" {
		t.Errorf("expected no violation for compliant KMS key, got: %s", res.Violation)
	}
}
