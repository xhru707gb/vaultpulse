// Package quarantine provides a thread-safe store for tracking secret paths
// that have been flagged for immediate attention — due to expiry, leakage,
// policy violation, or manual intervention.
//
// Usage:
//
//	store := quarantine.New()
//	_ = store.Add("secret/db/prod", quarantine.ReasonLeaked, "found in logs")
//	if store.IsQuarantined("secret/db/prod") { ... }
package quarantine
