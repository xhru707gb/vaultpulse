// Package audit provides structured audit logging and reporting for
// vaultpulse secret-check and rotation events.
//
// # Logger
//
// NewLogger returns a Logger that appends newline-delimited JSON records to any
// io.Writer. Each record captures the UTC timestamp, secret path, event type,
// outcome status, and remaining TTL.
//
// # Summary
//
// Summary reads a previously written audit log and returns per-path event
// counts grouped by status, useful for dashboards or CI gate checks.
//
// # Report
//
// Report reads a previously written audit log and renders a human-readable
// tabular report to any io.Writer, suitable for terminal output or file export.
package audit
