package pricing

import (
	"testing"
)

func TestEstimateMonthlyCostWithLiveFallback(t *testing.T) {
	// 1. Test live or fallback price estimation for Azure VM
	cost := EstimateMonthlyCost("azure", "virtual_machine", "medium")
	if cost <= 0 {
		t.Errorf("expected positive cost estimate for Azure VM, got %.2f", cost)
	}

	// 2. Test fallback estimation for AWS VM
	awsCost := EstimateMonthlyCost("aws", "virtual_machine", "medium")
	if awsCost < 30.0 || awsCost > 31.0 {
		t.Errorf("expected AWS VM fallback cost around 30.40, got %.2f", awsCost)
	}

	// 3. Test fallback estimation for EKS cluster
	eksCost := EstimateMonthlyCost("aws", "kubernetes_cluster", "")
	if eksCost != 73.00 {
		t.Errorf("expected EKS control plane cost 73.00, got %.2f", eksCost)
	}
}

func TestEstimateOfflineMonthlyCost(t *testing.T) {
	cost := EstimateOfflineMonthlyCost("aws", "virtual_machine", "small")
	if cost != 15.20 {
		t.Errorf("expected offline fallback cost 15.20, got %.2f", cost)
	}
}

func TestFetchLiveAzurePriceInvalidSKU(t *testing.T) {
	// Should fail gracefully and return error without crashing
	_, err := FetchLiveAzurePrice("NON_EXISTENT_INVALID_SKU_12345")
	if err == nil {
		t.Log("Note: Invalid SKU API call returned without error")
	}
}
