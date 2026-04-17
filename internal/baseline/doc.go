// Package baseline provides point-in-time recording of Vault secret metadata
// and drift detection by comparing current state against the recorded baseline.
//
// Typical usage:
//
//	store := baseline.New()
//	store.Record(baseline.Entry{Path: "secret/db", Version: 3, TTL: 3600})
//	drifts, err := store.Compare(current)
package baseline
