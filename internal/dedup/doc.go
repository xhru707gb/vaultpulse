// Package dedup implements alert deduplication for vaultpulse.
//
// A Deduplicator tracks (path, event) pairs and suppresses repeated
// notifications that occur within a configurable time window. This prevents
// alert fatigue when the same secret remains in an expiring or overdue state
// across multiple check cycles.
//
// Usage:
//
//	d, _ := dedup.New(15 * time.Minute)
//	if !d.IsDuplicate(path, "expired") {
//		// send alert
//	}
package dedup
