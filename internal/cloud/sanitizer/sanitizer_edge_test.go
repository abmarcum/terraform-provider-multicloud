package sanitizer

import (
	"testing"
)

func TestSanitizeResourceNamePathTraversal(t *testing.T) {
	// Path traversal attempt should strip '../' and directory delimiters
	raw := "../../etc/passwd"
	sanitized := SanitizeResourceName(raw, "aws", "storage_bucket")

	if sanitized == raw || sanitized == "../../etc/passwd" {
		t.Errorf("SanitizeResourceName failed to sanitize path traversal input: %s", sanitized)
	}
}

func TestSanitizeResourceNameShortLength(t *testing.T) {
	// Azure storage account names shorter than 3 chars should append fallback suffix
	raw := "ab"
	sanitized := SanitizeResourceName(raw, "azure", "storage_bucket")

	if len(sanitized) < 3 {
		t.Errorf("expected sanitized azure storage name to be >= 3 chars, got '%s'", sanitized)
	}
}

func TestSanitizeResourceNameNestedTraversalAndNonStorage(t *testing.T) {
	// Nested traversal bypass attempt ("....//") and query/fragment injection on non-storage resource
	raw := "....//....//admin?override=true#frag"
	sanitized := SanitizeResourceName(raw, "gcp", "virtual_machine")

	if sanitized != "adminoverridetruefrag" {
		t.Errorf("expected nested traversal and URL chars to be stripped, got %q", sanitized)
	}
}

