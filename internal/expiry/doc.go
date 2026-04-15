// Package expiry provides secret expiry checking and status reporting
// for the vaultpulse CLI.
//
// The Checker type queries HashiCorp Vault secret metadata and evaluates
// each secret's remaining TTL against a configurable warning threshold,
// classifying secrets as OK, WARNING, or EXPIRED.
//
// The FormatTable function renders a human-readable tabular summary of
// multiple SecretStatus values, with optional ANSI colour highlighting.
//
// Typical usage:
//
//	checker := expiry.NewChecker(vaultClient, 24*time.Hour)
//	statuses, err := checker.CheckAll(cfg.SecretPaths)
//	if err != nil {
//		log.Fatal(err)
//	}
//	expiry.FormatTable(os.Stdout, statuses, true)
package expiry
