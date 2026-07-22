package telemetry

import (
	"testing"
	"time"
)

func TestCacheEngineSetAndGet(t *testing.T) {
	cache := NewCacheEngine()

	cache.Set("azure-price-vm", 30.10, 100*time.Millisecond)

	val, found := cache.Get("azure-price-vm")
	if !found {
		t.Fatalf("expected key 'azure-price-vm' to be found in cache")
	}

	if price, ok := val.(float64); !ok || price != 30.10 {
		t.Errorf("expected cached price 30.10, got %v", val)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)
	_, found = cache.Get("azure-price-vm")
	if found {
		t.Errorf("expected expired key 'azure-price-vm' to return false")
	}
}
