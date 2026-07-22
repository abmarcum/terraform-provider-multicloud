package telemetry

import (
	"context"
	"sync"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/pricing"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/security"
)

// EvalTask defines a concurrent evaluation task
type EvalTask struct {
	ResourceName string
	ProviderType string
	ResourceType string
	Public       bool
	Encrypted    bool
}

// EvalResult holds findings from a concurrent evaluation task
type EvalResult struct {
	ResourceName string
	Findings     []security.SecurityAuditFinding
	CostEstimate float64
}

// ConcurrentEvaluator evaluates security audits and cost estimates in parallel
func ConcurrentEvaluator(ctx context.Context, tasks []EvalTask) []EvalResult {
	results := make([]EvalResult, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t EvalTask) {
			defer wg.Done()

			findings := security.AuditResource(t.ProviderType, t.ResourceType, t.ResourceName, t.Public, t.Encrypted)
			cost := pricing.EstimateMonthlyCost(t.ProviderType, t.ResourceType, "medium")

			results[idx] = EvalResult{
				ResourceName: t.ResourceName,
				Findings:     findings,
				CostEstimate: cost,
			}
		}(i, task)
	}

	wg.Wait()
	return results
}
