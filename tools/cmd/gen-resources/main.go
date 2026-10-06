package main

import (
	"context"
	"fmt"
	"os"

	"github.com/abmarcum/multi-cloud-provider/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func main() {
	fmt.Println("[gen-resources] Resource Code Generator & Schema Validator initialized.")

	ctx := context.Background()
	p := provider.New("1.0.0")()
	factories := p.Resources(ctx)

	fmt.Printf("[gen-resources] Validating resource definitions across %d unified schema targets.\n", len(factories))

	verifiedCount := 0
	totalAttributes := 0
	for _, factory := range factories {
		res := factory()

		metaReq := resource.MetadataRequest{ProviderTypeName: "multicloud"}
		metaResp := &resource.MetadataResponse{}
		res.Metadata(ctx, metaReq, metaResp)

		schemaReq := resource.SchemaRequest{}
		schemaResp := &resource.SchemaResponse{}
		res.Schema(ctx, schemaReq, schemaResp)

		if metaResp.TypeName == "" || schemaResp.Schema.Attributes == nil {
			fmt.Printf("[gen-resources] ERROR: Invalid metadata or nil schema attributes on resource %s\n", metaResp.TypeName)
			os.Exit(1)
		}

		for _, reqAttr := range []string{"id", "region", "extra_config"} {
			if _, ok := schemaResp.Schema.Attributes[reqAttr]; !ok {
				fmt.Printf("[gen-resources] ERROR: Resource %s missing mandatory '%s' attribute\n", metaResp.TypeName, reqAttr)
				os.Exit(1)
			}
		}

		totalAttributes += len(schemaResp.Schema.Attributes)
		verifiedCount++
	}

	fmt.Printf("[gen-resources] All %d resource code definitions verified (%d total schema attributes).\n", verifiedCount, totalAttributes)
}
