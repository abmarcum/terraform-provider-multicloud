package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abmarcum/multi-cloud-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func main() {
	fmt.Println("======================================================================")
	fmt.Println("  IDP DEVELOPER PORTAL EXPORTER (KRATIX PROMISES & BACKSTAGE TEMPLATES)")
	fmt.Println("======================================================================")

	ctx := context.Background()
	p := provider.New("1.0.0")()
	factories := p.Resources(ctx)

	resourceTypes := make([]string, 0, len(factories))
	for _, factory := range factories {
		res := factory()
		metaResp := &resource.MetadataResponse{}
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "multicloud"}, metaResp)
		resourceTypes = append(resourceTypes, metaResp.TypeName)
	}
	sort.Strings(resourceTypes)

	var enumLines strings.Builder
	for _, rt := range resourceTypes {
		enumLines.WriteString(fmt.Sprintf("                        - %s\n", rt))
	}

	kratixPromiseYAML := fmt.Sprintf(`apiVersion: platform.kratix.io/v1alpha1
kind: Promise
metadata:
  name: multicloud-infrastructure-promise
spec:
  api:
    apiVersion: apiextensions.k8s.io/v1
    kind: CustomResourceDefinition
    metadata:
      name: multicloudresources.marketplace.kratix.io
    spec:
      group: marketplace.kratix.io
      names:
        kind: MultiCloudResource
        plural: multicloudresources
        singular: multicloudresource
      scope: Namespaced
      versions:
        - name: v1alpha1
          served: true
          storage: true
          schema:
            openAPIV3Schema:
              type: object
              properties:
                spec:
                  type: object
                  required:
                    - resourceType
                    - providerType
                    - name
                  properties:
                    resourceType:
                      type: string
                      enum:
%s                    providerType:
                      type: string
                      enum:
                        - aws
                        - gcp
                        - azure
                    name:
                      type: string
                    region:
                      type: string
  workflows:
    resource:
      configure:
        - apiVersion: platform.kratix.io/v1alpha1
          kind: Pipeline
          metadata:
            name: terraform-multicloud-provision
          spec:
            containers:
              - name: terraform-apply-step
                image: hashicorp/terraform:1.9.8
`, enumLines.String())

	cwd, _ := os.Getwd()
	outPath := filepath.Clean(filepath.Join(cwd, "kratix_multicloud_promise.yaml"))
	/* #nosec G306 */
	if err := os.WriteFile(outPath, []byte(kratixPromiseYAML), 0600); err != nil {
		fmt.Printf("Error writing Kratix Promise YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[IDP Exporter] Exported Kratix Promise manifest (%d unified resource types) to %s\n", len(resourceTypes), outPath)
}
