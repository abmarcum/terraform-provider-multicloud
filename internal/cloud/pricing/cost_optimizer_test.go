package pricing

import (
	"testing"
)

func TestRecommendCostOptimizationsAWS(t *testing.T) {
	rec := RecommendCostOptimizations("aws", "my-app-server", "medium")
	if rec == nil {
		t.Fatalf("expected cost recommendation for medium AWS instance, got nil")
	}

	if rec.EstimatedSaving != 6.08 {
		t.Errorf("expected 6.08 savings, got %.2f", rec.EstimatedSaving)
	}

	if rec.SuggestedTier != "t4g.medium (AWS Graviton2)" {
		t.Errorf("unexpected suggested tier: %s", rec.SuggestedTier)
	}
}

func TestRecommendCostOptimizationsSmallTier(t *testing.T) {
	rec := RecommendCostOptimizations("gcp", "small-vm", "small")
	if rec != nil {
		t.Errorf("expected nil recommendation for small tier VM, got %v", rec)
	}
}
