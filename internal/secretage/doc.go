// Package secretage tracks the age of Vault secrets and surfaces
// those that have exceeded their configured maximum lifetime.
//
// Usage:
//
//	tracker := secretage.New()
//	_ = tracker.Register("secret/db", createdAt, 30*24*time.Hour)
//	statuses := tracker.EvaluateAll()
//	fmt.Print(secretage.FormatTable(statuses))
package secretage
