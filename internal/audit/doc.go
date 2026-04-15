// Package audit provides structured JSON audit logging for VaultPulse
// operations. Each audit entry records the event type, affected secret path,
// outcome, and a UTC timestamp so that operators can reconstruct a full
// history of expiry checks, rotation evaluations, and webhook dispatches.
//
// Usage:
//
//	logger, err := audit.NewLogger("/var/log/vaultpulse/audit.jsonl")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer logger.Close()
//
//	logger.Log(audit.Entry{
//		Event:   "expiry_check",
//		Path:    "secret/db/password",
//		Outcome: "warning",
//	})
package audit
