package telemetry

import (
	"sync"
	"time"
)

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// CacheEngine provides a thread-safe in-memory state memoization cache with TTL support
type CacheEngine struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

// NewCacheEngine initializes a new CacheEngine
func NewCacheEngine() *CacheEngine {
	return &CacheEngine{
		items: make(map[string]cacheItem),
	}
}

// Set stores a key-value pair with a specific Time-To-Live (TTL)
func (c *CacheEngine) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a key-value pair if not expired
func (c *CacheEngine) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.value, true
}

// Clear removes all cached items
func (c *CacheEngine) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}
