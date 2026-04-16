// Package retry implements exponential backoff retry logic for use
// throughout vaultpulse when making calls to external systems such as
// HashiCorp Vault or webhook endpoints.
//
// # Usage
//
//	cfg := retry.DefaultConfig()
//	err := retry.Do(ctx, cfg, func() error {
//		return client.Ping(ctx)
//	})
//
// The delay between attempts doubles on each retry, starting at
// BaseDelay and capped at MaxDelay. Context cancellation is respected
// between attempts so callers can abort early.
package retry
