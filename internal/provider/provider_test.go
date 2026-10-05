package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := New("0.1.0")()
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.Background(), req, resp)

	if resp.TypeName != "multicloud" {
		t.Errorf("expected TypeName 'multicloud', got '%s'", resp.TypeName)
	}
	if resp.Version != "0.1.0" {
		t.Errorf("expected Version '0.1.0', got '%s'", resp.Version)
	}
}

func TestProviderResourcesCount(t *testing.T) {
	p := &MulticloudProvider{version: "0.1.0"}
	resources := p.Resources(context.Background())

	if len(resources) != 70 {
		t.Errorf("expected 70 registered unified resources, got %d", len(resources))
	}
}

func TestProviderDataSourcesCount(t *testing.T) {
	p := &MulticloudProvider{version: "0.1.0"}
	dataSources := p.DataSources(context.Background())

	if len(dataSources) != 2 {
		t.Errorf("expected 2 registered data sources, got %d", len(dataSources))
	}
}
