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

func TestRecommendCostOptimizationsLargeAndMultiCloud(t *testing.T) {
	awsLarge := RecommendCostOptimizations("aws", "prod-db", "large")
	if awsLarge == nil || awsLarge.EstimatedSaving != 12.16 {
		t.Errorf("expected 12.16 savings for large AWS instance, got %+v", awsLarge)
	}

	gcpMed := RecommendCostOptimizations("gcp", "gcp-worker", "medium")
	if gcpMed == nil || gcpMed.EstimatedSaving != 5.96 {
		t.Errorf("expected 5.96 savings for medium GCP instance, got %+v", gcpMed)
	}

	gcpLarge := RecommendCostOptimizations("gcp", "gcp-analytics", "n2-standard-4")
	if gcpLarge == nil || gcpLarge.EstimatedSaving != 11.92 {
		t.Errorf("expected 11.92 savings for large GCP instance, got %+v", gcpLarge)
	}

	azMed := RecommendCostOptimizations("azure", "az-web", "medium")
	if azMed == nil || azMed.EstimatedSaving != 6.02 {
		t.Errorf("expected 6.02 savings for medium Azure instance, got %+v", azMed)
	}

	azLarge := RecommendCostOptimizations("azure", "az-batch", "large")
	if azLarge == nil || azLarge.EstimatedSaving != 12.04 {
		t.Errorf("expected 12.04 savings for large Azure instance, got %+v", azLarge)
	}
}

func TestRecommendCostOptimizationsSmallTier(t *testing.T) {
	rec := RecommendCostOptimizations("gcp", "small-vm", "small")
	if rec != nil {
		t.Errorf("expected nil recommendation for small tier VM, got %v", rec)
	}
}
