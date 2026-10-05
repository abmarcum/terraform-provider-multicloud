package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDriftDetectorMain(t *testing.T) {
	main()
}

func TestEvaluateStateDriftWithStateFile(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, "terraform.tfstate")
	stateContent := `{
  "resources": [
    {
      "mode": "managed",
      "type": "multicloud_storage_bucket",
      "name": "compliant_bucket",
      "instances": [
        {
          "attributes": {
            "provider_type": "aws",
            "region": "us-east-1",
            "encryption_enabled": true
          }
        }
      ]
    },
    {
      "mode": "managed",
      "type": "multicloud_storage_bucket",
      "name": "unencrypted_bucket",
      "instances": [
        {
          "attributes": {
            "provider_type": "gcp",
            "region": "us-central1",
            "encryption_enabled": false
          }
        }
      ]
    }
  ]
}`
	if err := os.WriteFile(stateFile, []byte(stateContent), 0600); err != nil {
		t.Fatalf("failed to write test tfstate: %v", err)
	}

	reports := evaluateStateDrift(context.Background(), stateFile)
	if len(reports) != 2 {
		t.Fatalf("expected 2 drift reports, got %d", len(reports))
	}
	if reports[0].HasDrift {
		t.Errorf("expected compliant_bucket to have HasDrift=false, got true (%s)", reports[0].DriftDetails)
	}
	if !reports[1].HasDrift {
		t.Errorf("expected unencrypted_bucket to have HasDrift=true")
	}
}
