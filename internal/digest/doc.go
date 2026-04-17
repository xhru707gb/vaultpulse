// Package digest aggregates expiry, rotation and health signals for a set of
// Vault secrets into a single periodic summary report.
//
// Usage:
//
//	b := digest.NewBuilder(nil)
//	report := b.Build(entries)
//	fmt.Print(digest.FormatTable(report))
package digest
