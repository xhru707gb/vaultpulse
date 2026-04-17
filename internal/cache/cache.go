// Package cache provides an in-memory TTL cache for Vault secret metadata
// to reduce repeated API calls during a single vaultpulse run.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached value along with its expiration time.
type Entry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// Cache is a simple thread-safe in-memory TTL store.
type Cache struct {
	mu         sync.RWMutex
	items      map[string]Entry
	defaultTTL time.Duration
}

// New creates a Cache with the given default TTL for all entries.
func New(defaultTTL time.Duration) *Cache {
	return &Cache{
		items:      make(map[string]Entry),
		defaultTTL: defaultTTL,
	}
}

// Set stores a value under key with the default TTL.
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value under key with an explicit TTL.
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = Entry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves a value by key. Returns (value, true) if present and not
// expired, or (nil, false) otherwise.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.items[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Value, true
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Entry)
}

// Len returns the number of entries currently in the cache (including expired).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Purge removes all expired entries from the cache and returns the number
// of entries that were deleted.
func (c *Cache) Purge() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	count := 0
	for key, entry := range c.items {
		if now.After(entry.ExpiresAt) {
			delete(c.items, key)
			count++
		}
	}
	return count
}
