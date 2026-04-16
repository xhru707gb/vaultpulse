// export_test.go exposes internal fields for white-box testing.
package cache

import "time"

// EntryExpiresAt returns the expiration time of a cache entry by key.
// It is only available during testing.
func (c *Cache) EntryExpiresAt(key string) (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.items[key]
	if !ok {
		return time.Time{}, false
	}
	return e.ExpiresAt, true
}

// Len returns the number of items currently in the cache.
// It is only available during testing.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
