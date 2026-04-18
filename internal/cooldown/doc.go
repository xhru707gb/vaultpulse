// Package cooldown provides a per-key cooldown tracker used to suppress
// repeated alerts or actions within a configurable quiet window.
//
// Usage:
//
//	tr, err := cooldown.New(5 * time.Minute)
//	if err != nil { ... }
//
//	if !tr.IsCoolingDown(secretPath) {
//		// send alert
//		tr.Record(secretPath)
//	}
package cooldown
