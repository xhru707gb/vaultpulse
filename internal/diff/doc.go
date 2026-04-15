// Package diff provides utilities for computing and displaying differences
// between two snapshots of Vault secret metadata.
//
// A diff identifies secrets that have been added, removed, or modified
// (e.g. version bump or expiry change) between two point-in-time snapshots
// captured by the snapshot package.
//
// Typical usage:
//
//	prev, _ := snapshot.Load("vault-snap-prev.json")
//	curr, _ := snapshot.Load("vault-snap-curr.json")
//	changes := diff.Compute(prev.Entries, curr.Entries)
//	fmt.Print(diff.FormatTable(changes))
package diff
