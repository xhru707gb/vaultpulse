// Package cache provides a lightweight, thread-safe in-memory TTL cache
// used by vaultpulse to store Vault secret metadata between repeated
// lookups within a single CLI invocation.
//
// Entries are keyed by the Vault secret path and expire after a
// configurable duration, preventing stale data from persisting across
// long-running watch loops.
//
// Usage:
//
//	c := cache.New(30 * time.Second)
//	c.Set("secret/myapp/db", metadata)
//	if v, ok := c.Get("secret/myapp/db"); ok {
//	    // use cached metadata
//	}
package cache
