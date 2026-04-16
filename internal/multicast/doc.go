// Package multicast provides a thread-safe broadcaster that fans out vault
// lifecycle events (expiry warnings, rotation alerts, health changes) to an
// arbitrary number of named handlers concurrently.
//
// Usage:
//
//	b := multicast.New()
//	b.Register("slack", slackHandler)
//	b.Register("audit", auditHandler)
//	b.Broadcast("expiry.warning", payload)
package multicast
