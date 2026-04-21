// Package secretlookup implements a reverse-lookup index that maps secret
// fingerprints to the vault paths that reference them.
//
// It enables vaultpulse to detect when the same secret value (identified by
// its SHA-256 fingerprint) is stored under multiple paths — a common
// security hygiene issue that increases blast radius on compromise.
//
// Usage:
//
//	idx := secretlookup.New()
//	_ = idx.Add("secret/db/prod", fingerprint)
//	_ = idx.Add("secret/db/staging", fingerprint)
//	dupes := idx.Duplicates() // returns both paths grouped by fingerprint
package secretlookup
