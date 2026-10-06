package pricing

import (
	"fmt"
	"strings"
)

// CostOptimizationRecommendation represents architecture savings recommendations
type CostOptimizationRecommendation struct {
	ResourceName    string
	CurrentProvider string
	CurrentTier     string
	SuggestedTier   string
	EstimatedSaving float64 // Monthly USD savings
	Message         string
}

// RecommendCostOptimizations analyzes instance tier or SKU selections and recommends Arm-based equivalent architectures
func RecommendCostOptimizations(providerType string, resourceName string, currentTier string) *CostOptimizationRecommendation {
	p := strings.ToLower(strings.TrimSpace(providerType))
	tier := strings.ToLower(strings.TrimSpace(currentTier))

	switch tier {
	case "medium", "t3.medium", "m6i.large", "e2-medium", "n2-standard-2", "standard_b2s", "standard_d2s_v5":
		switch p {
		case "aws":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "aws",
				CurrentTier:     currentTier,
				SuggestedTier:   "t4g.medium (AWS Graviton2)",
				EstimatedSaving: 6.08, // 20% savings
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to AWS Graviton2 (t4g.medium) for 20%% cost savings ($6.08 USD/mo).", resourceName),
			}
		case "gcp":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "gcp",
				CurrentTier:     currentTier,
				SuggestedTier:   "t2a-standard-2 (GCP Tau ARM)",
				EstimatedSaving: 5.96, // 20% savings
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to GCP Tau ARM (t2a-standard-2) for 20%% cost savings ($5.96 USD/mo).", resourceName),
			}
		case "azure":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "azure",
				CurrentTier:     currentTier,
				SuggestedTier:   "Standard_D2ps_v5 (Ampere Altra)",
				EstimatedSaving: 6.02, // 20% savings
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to Azure Ampere Altra for 20%% cost savings ($6.02 USD/mo).", resourceName),
			}
		}

	case "large", "t3.large", "m6i.xlarge", "e2-standard-2", "n2-standard-4", "standard_b2ms", "standard_d4s_v5":
		switch p {
		case "aws":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "aws",
				CurrentTier:     currentTier,
				SuggestedTier:   "t4g.large (AWS Graviton2)",
				EstimatedSaving: 12.16, // 20% savings on $60.80/mo
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to AWS Graviton2 (t4g.large) for 20%% cost savings ($12.16 USD/mo).", resourceName),
			}
		case "gcp":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "gcp",
				CurrentTier:     currentTier,
				SuggestedTier:   "t2a-standard-4 (GCP Tau ARM)",
				EstimatedSaving: 11.92, // 20% savings on $59.60/mo
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to GCP Tau ARM (t2a-standard-4) for 20%% cost savings ($11.92 USD/mo).", resourceName),
			}
		case "azure":
			return &CostOptimizationRecommendation{
				ResourceName:    resourceName,
				CurrentProvider: "azure",
				CurrentTier:     currentTier,
				SuggestedTier:   "Standard_D4ps_v5 (Ampere Altra)",
				EstimatedSaving: 12.04, // 20% savings on $60.20/mo
				Message:         fmt.Sprintf("[CostOptimizer] Resource '%s': Switch to Azure Ampere Altra (Standard_D4ps_v5) for 20%% cost savings ($12.04 USD/mo).", resourceName),
			}
		}
	}

	return nil
}
