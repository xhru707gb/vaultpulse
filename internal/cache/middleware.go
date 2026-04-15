package cache

import (
	"fmt"
	"time"
)

// Fetcher is a function that retrieves a value from an upstream source.
type Fetcher[T any] func(key string) (T, error)

// GetOrFetch returns the cached value for key if present and valid.
// On a cache miss it calls fetch, stores the result, and returns it.
// If fetch returns an error the result is not cached.
func GetOrFetch[T any](c *Cache, key string, ttl time.Duration, fetch Fetcher[T]) (T, error) {
	if raw, ok := c.Get(key); ok {
		if typed, ok := raw.(T); ok {
			return typed, nil
		}
		// Stale type in cache — evict and re-fetch.
		c.Delete(key)
	}

	value, err := fetch(key)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("cache: fetch failed for key %q: %w", key, err)
	}

	if ttl > 0 {
		c.SetWithTTL(key, value, ttl)
	} else {
		c.Set(key, value)
	}

	return value, nil
}

// Invalidate removes all keys that match the given prefix.
func Invalidate(c *Cache, prefix string) int {
	keys := c.Keys()
	removed := 0
	for _, k := range keys {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			c.Delete(k)
			removed++
		}
	}
	return removed
}
