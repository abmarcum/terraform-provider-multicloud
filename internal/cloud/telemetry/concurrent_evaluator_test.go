package telemetry

import (
	"context"
	"testing"
)

func TestConcurrentEvaluator(t *testing.T) {
	tasks := []EvalTask{
		{ResourceName: "bucket-aws", ProviderType: "aws", ResourceType: "storage_bucket", Public: false, Encrypted: true},
		{ResourceName: "vm-gcp", ProviderType: "gcp", ResourceType: "virtual_machine", Public: true, Encrypted: true},
		{ResourceName: "container-azure", ProviderType: "azure", ResourceType: "storage_bucket", Public: true, Encrypted: false},
	}

	results := ConcurrentEvaluator(context.Background(), tasks)

	if len(results) != 3 {
		t.Fatalf("expected 3 concurrent evaluation results, got %d", len(results))
	}

	if results[0].ResourceName != "bucket-aws" || results[1].ResourceName != "vm-gcp" {
		t.Errorf("unexpected task result ordering in ConcurrentEvaluator")
	}

	if len(results[2].Findings) != 2 {
		t.Errorf("expected 2 findings for unencrypted public azure container, got %d", len(results[2].Findings))
	}
}
