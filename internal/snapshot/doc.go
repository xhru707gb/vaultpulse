// Package snapshot captures and persists a point-in-time view of Vault secret
// expiry and rotation statuses.
//
// A Snapshot encodes both expiry and rotation states alongside a UTC timestamp,
// enabling offline diffing, historical trending, and audit trail generation
// without requiring a live Vault connection.
//
// Usage:
//
//	w := snapshot.NewWriter("/var/lib/vaultpulse/latest.json")
//	if err := w.Write(expiryStatuses, rotationStatuses); err != nil {
//		log.Fatal(err)
//	}
//
//	snap, err := snapshot.Load("/var/lib/vaultpulse/latest.json")
package snapshot
