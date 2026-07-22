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
}
