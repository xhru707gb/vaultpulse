// Package ratelimit implements a token-bucket rate limiter used to throttle
// outbound requests to the HashiCorp Vault API.
//
// Usage:
//
//	limiter, err := ratelimit.New(ratelimit.Config{
//		RequestsPerSecond: 10,
//		Burst:             5,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Before each Vault API call:
//	limiter.Wait()
//
// The limiter is safe for concurrent use.
package ratelimit
