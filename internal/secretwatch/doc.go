// Package secretwatch implements a lightweight periodic watcher for HashiCorp
// Vault secret paths. It compares successive snapshots of secret metadata
// returned by a user-supplied SecretLister and fires a Handler with the
// detected changes (added, removed, or modified secrets).
//
// Usage:
//
//	w, err := secretwatch.New(myLister, myHandler, 30*time.Second)
//	if err != nil { ... }
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	w.Run(ctx)
package secretwatch
