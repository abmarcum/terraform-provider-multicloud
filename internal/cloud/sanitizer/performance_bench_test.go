package sanitizer

import (
	"testing"
	"time"

	"github.com/abmarcum/multi-cloud-provider/internal/cloud/pricing"
	"github.com/abmarcum/multi-cloud-provider/internal/cloud/telemetry"
)

func BenchmarkSanitizeResourceName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = SanitizeResourceName("My_Long_Azure_Storage_Account_Name_2026", "azure", "storage_bucket")
	}
}

func BenchmarkEstimateMonthlyCost(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = pricing.EstimateMonthlyCost("aws", "virtual_machine", "medium")
	}
}

func BenchmarkCacheEngine(b *testing.B) {
	cache := telemetry.NewCacheEngine()
	cache.Set("test-key", "test-value", 10*time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cache.Get("test-key")
	}
}
