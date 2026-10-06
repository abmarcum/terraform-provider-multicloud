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
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= 1024 {
		for k, v := range c.items {
			if now.After(v.expiresAt) {
				delete(c.items, k)
			}
		}
	}

	c.items[key] = cacheItem{
		value:     value,
		expiresAt: now.Add(ttl),
	}
}

// Get retrieves a key-value pair if not expired, evicting expired entries lazily
func (c *CacheEngine) Get(key string) (interface{}, bool) {
	now := time.Now()
	c.mu.RLock()
	item, found := c.items[key]
	if !found {
		c.mu.RUnlock()
		return nil, false
	}
	if !now.After(item.expiresAt) {
		val := item.value
		c.mu.RUnlock()
		return val, true
	}
	c.mu.RUnlock()

	// Evict expired item under write lock
	c.mu.Lock()
	if current, ok := c.items[key]; ok && now.After(current.expiresAt) {
		delete(c.items, key)
	}
	c.mu.Unlock()
	return nil, false
}

// Clear removes all cached items
func (c *CacheEngine) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}
