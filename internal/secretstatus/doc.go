// Package secretstatus aggregates health status across multiple secret
// evaluation dimensions (expiry, rotation, policy compliance, risk score).
//
// Callers register Provider implementations — each representing one
// dimension — and call Evaluate or EvaluateAll to obtain a combined
// Level (OK / Warning / Critical) together with human-readable reasons.
//
// Example usage:
//
//	eval, err := secretstatus.New(expiryProvider, rotationProvider)
//	if err != nil { ... }
//	entry, err := eval.Evaluate("secret/my-app/db")
package secretstatus
