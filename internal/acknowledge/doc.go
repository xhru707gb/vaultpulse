// Package acknowledge provides alert acknowledgement tracking for vaultpulse.
//
// An acknowledged secret path will be suppressed from alerting until the
// acknowledgement window expires or is explicitly revoked. This allows
// operators to silence known issues without disabling monitoring globally.
//
// Usage:
//
//	tr, err := acknowledge.New(4 * time.Hour)
//	_ = tr.Acknowledge("secret/db/password", "operator")
//	if tr.IsAcknowledged("secret/db/password") {
//		// skip alerting
//	}
package acknowledge
