// Package secretgroup provides grouping of Vault secret paths into named
// collections for bulk policy evaluation, reporting, and alerting.
//
// Usage:
//
//	g := secretgroup.New()
//	_ = g.Add("databases", "secret/db/prod")
//	_ = g.Add("databases", "secret/db/staging")
//	groups := g.All()
//	fmt.Print(secretgroup.FormatTable(groups))
package secretgroup
