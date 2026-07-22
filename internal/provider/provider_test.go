package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := New("1.0.0")()
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.Background(), req, resp)

	if resp.TypeName != "multicloud" {
		t.Errorf("expected TypeName 'multicloud', got '%s'", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected Version '1.0.0', got '%s'", resp.Version)
	}
}

func TestProviderResourcesCount(t *testing.T) {
	p := &MulticloudProvider{version: "1.0.0"}
	resources := p.Resources(context.Background())

	if len(resources) != 33 {
		t.Errorf("expected 33 registered unified resources, got %d", len(resources))
	}
}
