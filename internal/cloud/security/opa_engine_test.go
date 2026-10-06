package security

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateOPARegoPolicy(t *testing.T) {
	// 1. Test passing OPA policy
	res := EvaluateOPARegoPolicy("aws", "multicloud_storage_bucket", "my-bucket", "must_have_tags")
	if !res.Passed {
		t.Errorf("expected OPA policy evaluation to pass, got violation: %s", res.Violation)
	}

	// 2. Test failing OPA policy
	res = EvaluateOPARegoPolicy("gcp", "multicloud_virtual_machine", "my-vm", "disallow_public_ip")
	if res.Passed {
		t.Errorf("expected OPA policy evaluation to fail for public VM, but it passed")
	}

	// 3. Test attribute-aware OPA policy evaluation
	res = EvaluateOPARegoPolicyWithAttrs("aws", "multicloud_storage_bucket", "my-bucket", "must_have_tags", map[string]interface{}{
		"tags": map[string]string{},
	})
	if res.Passed {
		t.Errorf("expected must_have_tags to fail when tags map is empty")
	}

	res = EvaluateOPARegoPolicyWithAttrs("gcp", "multicloud_virtual_machine", "my-vm", "disallow_public_ip", map[string]interface{}{
		"associate_public_ip": false,
	})
	if !res.Passed {
		t.Errorf("expected disallow_public_ip to pass when associate_public_ip is false")
	}

	// 4. Test require_encryption and require_multi_az
	res = EvaluateOPARegoPolicyWithAttrs("aws", "multicloud_storage_bucket", "unenc-bucket", "require_encryption", map[string]interface{}{
		"encryption_enabled": false,
	})
	if res.Passed {
		t.Errorf("expected require_encryption to fail when encryption_enabled is false")
	}

	res = EvaluateOPARegoPolicyWithAttrs("aws", "multicloud_db_instance", "single-az-db", "require_multi_az", map[string]interface{}{
		"multi_az": false,
	})
	if res.Passed {
		t.Errorf("expected require_multi_az to fail when multi_az is false")
	}
}

func TestEvaluateRegoModuleAndPolicyPath(t *testing.T) {
	regoSrc := `
package multicloud.compliance

deny[msg] {
    input.resource_type == "multicloud_storage_bucket"
    input.attributes.encryption_enabled == false
    msg := "Custom Rego policy: storage buckets must enable encryption"
}
`
	res := EvaluateRegoModule("aws", "multicloud_storage_bucket", "test-bkt", regoSrc, map[string]interface{}{
		"encryption_enabled": false,
	})
	if res.Passed {
		t.Fatalf("expected custom Rego deny rule to trigger violation")
	}

	tmpDir := t.TempDir()
	policyFile := filepath.Join(tmpDir, "policy.rego")
	if err := os.WriteFile(policyFile, []byte(regoSrc), 0600); err != nil {
		t.Fatalf("failed to write temp rego file: %v", err)
	}
	t.Setenv("OPA_POLICY_PATH", policyFile)

	res = EvaluateOPARegoPolicyWithAttrs("aws", "multicloud_storage_bucket", "test-bkt", "custom_rule", map[string]interface{}{
		"encryption_enabled": false,
	})
	if res.Passed {
		t.Errorf("expected OPA_POLICY_PATH evaluation to fail on unencrypted bucket")
	}

	t.Setenv("OPA_POLICY_PATH", filepath.Join(tmpDir, "nonexistent.rego"))
	res = EvaluateOPARegoPolicyWithAttrs("aws", "multicloud_storage_bucket", "test-bkt", "custom_rule", map[string]interface{}{
		"encryption_enabled": true,
	})
	if res.Passed {
		t.Errorf("expected missing OPA_POLICY_PATH file to fail closed")
	}
}
