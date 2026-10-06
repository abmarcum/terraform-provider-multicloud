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

func TestSanitizeCloudIdentifierAndStripControlChars(t *testing.T) {
	injectedRegion := "us-east-1@attacker.com:8443/path"
	cleanRegion := SanitizeCloudIdentifier(injectedRegion, "us-east-1")
	if cleanRegion != "us-east-1attackercom8443path" {
		t.Errorf("expected URL authority characters to be stripped from region, got %q", cleanRegion)
	}

	if fallback := SanitizeCloudIdentifier("../../../", "us-central1"); fallback != "us-central1" {
		t.Errorf("expected fallback 'us-central1', got %q", fallback)
	}

	ansiInput := "prod-bucket\x1b[31m-spoofed\x07"
	if stripped := StripControlChars(ansiInput); stripped != "prod-bucket[31m-spoofed" {
		t.Errorf("expected control characters to be stripped, got %q", stripped)
	}
}

