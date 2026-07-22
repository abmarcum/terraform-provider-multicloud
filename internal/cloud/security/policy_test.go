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
}
