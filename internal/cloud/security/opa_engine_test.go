package security

import (
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
}

