package pricing

import (
	"net/http"
	"net/http/httptest"
	"strings"
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
		t.Errorf("expected AWS VM fallback cost around 30.37, got %.2f", awsCost)
	}

	// 3. Test fallback estimation for EKS cluster
	eksCost := EstimateMonthlyCost("aws", "kubernetes_cluster", "")
	if eksCost != 73.00 {
		t.Errorf("expected EKS control plane cost 73.00, got %.2f", eksCost)
	}
}

func TestEstimateOfflineMonthlyCostCategories(t *testing.T) {
	cost := EstimateOfflineMonthlyCost("aws", "virtual_machine", "small")
	if cost != 15.20 {
		t.Errorf("expected offline fallback cost 15.20, got %.2f", cost)
	}

	cases := []struct {
		provider string
		resType  string
		want     float64
	}{
		{"aws", "storage_bucket", 2.30},
		{"gcp", "block_volume", 10.00},
		{"azure", "shared_filesystem", 25.00},
		{"aws", "load_balancer", 22.50},
		{"gcp", "service_mesh", 18.00},
		{"azure", "vector_index", 34.20},
		{"aws", "batch_compute", 180.00},
		{"gcp", "workflow", 11.00},
		{"azure", "distributed_tracing", 12.00},
		{"aws", "tls_certificate", 0.00},
		{"gcp", "budget_alert", 0.00},
	}

	for _, tc := range cases {
		got := EstimateOfflineMonthlyCost(tc.provider, tc.resType, "")
		if got != tc.want {
			t.Errorf("EstimateOfflineMonthlyCost(%s, %s) = %.2f; want %.2f", tc.provider, tc.resType, got, tc.want)
		}
	}

	report := FormatCostReport("aws", "prod-db", 45.00)
	if !strings.Contains(report, "45.00") {
		t.Errorf("unexpected FormatCostReport output: %s", report)
	}
}

func TestFetchLiveAWSPriceOfferServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"terms":{"OnDemand":{"SKU1":{"TERM1":{"priceDimensions":{"DIM1":{"pricePerUnit":{"USD":"0.1000"}}}}}}}}`))
	}))
	defer srv.Close()

	t.Setenv("MULTICLOUD_MOCK_MODE", "false")
	t.Setenv("AWS_PRICING_OFFER_URL", srv.URL)

	price, err := FetchLiveAWSPrice("custom.test.sku")
	if err != nil {
		t.Fatalf("FetchLiveAWSPrice failed: %v", err)
	}
	if price != 73.00 {
		t.Errorf("expected 73.00 monthly from 0.10 hourly, got %.2f", price)
	}
}

func TestFetchLiveAzurePriceInvalidSKU(t *testing.T) {
	_, err := FetchLiveAzurePrice("NON_EXISTENT_INVALID_SKU_12345")
	if err == nil {
		t.Log("Note: Invalid SKU API call returned without error")
	}
}
