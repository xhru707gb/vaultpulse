// Package diff compares two snapshots of Vault secret metadata and reports
// which secrets have been added, removed, or modified between runs.
//
// Usage:
//
//	prev := []diff.SecretEntry{ /* loaded from previous snapshot */ }
//	curr := []diff.SecretEntry{ /* loaded from current check */ }
//	changes := diff.Compute(prev, curr)
//	for _, c := range changes {
//		fmt.Println(c)
//	}
//
// Changes are classified as ChangeAdded, ChangeRemoved, or ChangeModified.
// A modification is triggered when either the version number or expiry timestamp
// differs between the two snapshots.
package diff
