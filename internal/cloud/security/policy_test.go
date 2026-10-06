package security

import (
	"testing"
)

func TestValidatePolicy(t *testing.T) {
	attrs := map[string]interface{}{
		"is_public": true,
	}

	violations := ValidatePolicy("multicloud_storage_bucket", "my-public-bucket", attrs)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for public storage bucket, got %d", len(violations))
	}

	if violations[0].RuleName != "NO_PUBLIC_STORAGE" {
		t.Errorf("expected RuleName 'NO_PUBLIC_STORAGE', got '%s'", violations[0].RuleName)
	}

	// Test STORAGE_ENCRYPTION_REQUIRED
	encViolations := ValidatePolicy("multicloud_storage_bucket", "unencrypted-bkt", map[string]interface{}{
		"encryption_enabled": false,
	})
	if len(encViolations) != 1 || encViolations[0].RuleName != "STORAGE_ENCRYPTION_REQUIRED" {
		t.Errorf("expected STORAGE_ENCRYPTION_REQUIRED violation, got %+v", encViolations)
	}

	// Test DB_PUBLIC_ACCESS_DISALLOWED
	dbViolations := ValidatePolicy("multicloud_db_instance", "public-db", map[string]interface{}{
		"publicly_accessible": true,
	})
	if len(dbViolations) != 1 || dbViolations[0].RuleName != "DB_PUBLIC_ACCESS_DISALLOWED" {
		t.Errorf("expected DB_PUBLIC_ACCESS_DISALLOWED violation, got %+v", dbViolations)
	}

	// Test KMS_ROTATION_REQUIRED
	kmsViolations := ValidatePolicy("multicloud_kms_key", "no-rotate-key", map[string]interface{}{
		"rotation_enabled": false,
	})
	if len(kmsViolations) != 1 || kmsViolations[0].RuleName != "KMS_ROTATION_REQUIRED" {
		t.Errorf("expected KMS_ROTATION_REQUIRED violation, got %+v", kmsViolations)
	}
}
