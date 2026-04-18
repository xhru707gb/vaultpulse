// Package metrics provides a lightweight, thread-safe collector for
// aggregated VaultPulse check results.
//
// A Collector accumulates counts from expiry, rotation, health, and
// policy subsystems and exposes a point-in-time Snapshot that can be
// rendered as a table or consumed programmatically.
//
// # Collector
//
// NewCollector returns a ready-to-use Collector. Record replaces the
// current Snapshot atomically; Get returns a copy of the latest values.
//
// # Snapshot
//
// Snapshot is a plain struct whose fields map 1-to-1 to the columns
// shown in the dashboard table. Zero values are valid and mean "no data
// collected yet".
//
// # Formatting
//
// FormatTable renders a Snapshot as a human-readable text table suitable
// for terminal output. For machine consumption, callers can access the
// Snapshot fields directly.
//
// Usage:
//
//	c := metrics.NewCollector()
//	c.Record(metrics.Snapshot{
//		TotalSecrets: 10,
//		Expired:      2,
//		Warning:      3,
//	})
//	fmt.Print(metrics.FormatTable(c.Get()))
package metrics
