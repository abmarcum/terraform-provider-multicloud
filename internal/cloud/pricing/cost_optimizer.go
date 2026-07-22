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

// RecommendCostOptimizations analyzes instance tier selections and recommends Arm-based or spot equivalent architectures
func RecommendCostOptimizations(providerType string, resourceName string, currentTier string) *CostOptimizationRecommendation {
	p := strings.ToLower(providerType)
	tier := strings.ToLower(currentTier)

	if tier == "medium" || tier == "large" {
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
	}

	return nil
}
