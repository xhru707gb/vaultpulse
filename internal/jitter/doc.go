// Package jitter provides randomised jitter helpers for duration-based
// operations such as retry back-off and polling intervals.
//
// Use [New] to create a Jitter with a configurable factor, then call
// [Jitter.Apply] to add a one-sided random delta or [Jitter.ApplyRange]
// to spread the duration symmetrically around the base value.
package jitter
