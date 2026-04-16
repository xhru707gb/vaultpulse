// Package trend provides time-series analysis of VaultPulse audit log events.
//
// It groups audit entries by secret path and event type, buckets them into
// configurable time windows, and exposes both raw Report structs and a
// human-readable ASCII table with spark-bar visualisation.
//
// Typical usage:
//
//	analyzer, _ := trend.NewAnalyzer(24 * time.Hour)
//	reports := analyzer.Analyse(entries)
//	fmt.Print(trend.FormatTable(reports))
package trend
