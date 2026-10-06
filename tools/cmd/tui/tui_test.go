package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTUIMain(t *testing.T) {
	main()
}

func TestLoadResourcesFromHCL(t *testing.T) {
	tmpDir := t.TempDir()
	hclPath := filepath.Join(tmpDir, "main.tf")
	hclContent := `
resource "multicloud_virtual_machine" "web" {
  provider_type       = "aws"
  vm_name             = "prod-web-vm"
  size_tier           = "medium"
  associate_public_ip = false
  encryption_enabled  = true
}
`
	if err := os.WriteFile(hclPath, []byte(hclContent), 0600); err != nil {
		t.Fatalf("failed to write test HCL file: %v", err)
	}

	resources := loadResourcesFromHCL(tmpDir)
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource from HCL, got %d", len(resources))
	}
	if resources[0].Name != "prod-web-vm" || resources[0].Provider != "aws" || resources[0].Tier != "medium" {
		t.Errorf("unexpected parsed HCL resource: %+v", resources[0])
	}
}
