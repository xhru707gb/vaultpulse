// Package metrics provides a lightweight, thread-safe collector for
// aggregated VaultPulse check results.
//
// A Collector accumulates counts from expiry, rotation, health, and
// policy subsystems and exposes a point-in-time Snapshot that can be
// rendered as a table or consumed programmatically.
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
