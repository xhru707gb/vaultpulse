// Package renew implements automatic secret and lease renewal tracking for
// VaultPulse. It monitors registered secret paths and signals when renewal
// should occur based on a configurable fraction of the lease TTL, ensuring
// credentials are refreshed before they expire.
package renew
