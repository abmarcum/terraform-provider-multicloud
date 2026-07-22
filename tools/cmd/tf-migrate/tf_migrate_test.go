package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTFMigrateHCLConversionWithExtraConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tf-migrate-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sampleLegacyHCL := `
resource "aws_s3_bucket" "prod_storage" {
  bucket_name = "my-aws-bucket"
  force_destroy = "true"
}

resource "google_compute_instance" "app_vm" {
  name = "my-gcp-vm"
  custom_gpu_type = "nvidia-tesla-t4"
}
`

	legacyFile := filepath.Join(tmpDir, "legacy.tf")
	if err := os.WriteFile(legacyFile, []byte(sampleLegacyHCL), 0600); err != nil {
		t.Fatalf("failed to write legacy HCL test file: %v", err)
	}

	matches := resourceRegex.FindAllStringSubmatch(sampleLegacyHCL, -1)
	if len(matches) != 2 {
		t.Fatalf("expected 2 resource matches, got %d", len(matches))
	}

	if !strings.Contains(matches[0][3], "force_destroy") {
		t.Errorf("expected force_destroy attribute in body match")
	}
}

func TestTFMigrateResourceMappings(t *testing.T) {
	foundStorage := false
	for _, m := range resourceMappings {
		if m.LegacyType == "aws_s3_bucket" && m.UnifiedType == "multicloud_storage_bucket" {
			foundStorage = true
			break
		}
	}

	if !foundStorage {
		t.Errorf("expected aws_s3_bucket mapping to multicloud_storage_bucket")
	}
}

func TestTFMigrateMain(t *testing.T) {
	os.Args = []string{"tf-migrate", "--dry-run"}
	main()
}
