// Package semaphore implements a simple counting semaphore used to cap
// the number of concurrent requests made against a HashiCorp Vault cluster.
//
// Usage:
//
//	sem, err := semaphore.New(10)
//	if err != nil { ... }
//	if err := sem.Acquire(ctx); err != nil { ... }
//	defer sem.Release()
package semaphore
