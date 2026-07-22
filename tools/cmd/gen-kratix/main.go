package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("======================================================================")
	fmt.Println("  IDP DEVELOPER PORTAL EXPORTER (KRATIX PROMISES & BACKSTAGE TEMPLATES)")
	fmt.Println("======================================================================")

	kratixPromiseYAML := `apiVersion: platform.kratix.io/v1alpha1
kind: Promise
metadata:
  name: multicloud-infrastructure-promise
spec:
  api:
    apiVersion: marketplace.kratix.io/v1alpha1
    kind: MultiCloudBucket
  workflows:
    resource:
      configure:
        - name: terraform-apply-step
          image: hashicorp/terraform:latest
`

	cwd, _ := os.Getwd()
	outPath := filepath.Join(cwd, "kratix_multicloud_promise.yaml")
	/* #nosec G306 */
	if err := os.WriteFile(outPath, []byte(kratixPromiseYAML), 0600); err != nil {
		fmt.Printf("Error writing Kratix Promise YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[IDP Exporter] Exported Kratix Promise manifest to %s\n", outPath)
}
