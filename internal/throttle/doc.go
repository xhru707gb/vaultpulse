// Package throttle provides a concurrency limiter for outbound Vault API
// requests. It uses a buffered channel as a semaphore so that at most N
// requests are in-flight simultaneously, preventing Vault from being
// overwhelmed during bulk secret checks or rotation evaluations.
//
// Usage:
//
//	t, _ := throttle.New(10)
//	err := t.Do(ctx, func() error {
//		return vaultClient.GetSecretMeta(path)
//	})
package throttle
