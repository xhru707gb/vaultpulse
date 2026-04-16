// Package backoff implements exponential back-off with optional full-jitter,
// suitable for retrying Vault API calls or webhook deliveries.
//
// Usage:
//
//	b, err := backoff.New(backoff.DefaultConfig())
//	if err != nil { ... }
//	for {
//		if err := doWork(); err == nil { break }
//		time.Sleep(b.Next())
//	}
package backoff
