// Package expiry provides secret expiration checking, formatting,
// and alerting capabilities for the vaultpulse CLI.
//
// It exposes:
//   - Checker: queries Vault secret metadata and evaluates TTL
//     against configurable warning and expiry thresholds, returning
//     a slice of Status values.
//   - FormatTable: renders expiry statuses as a human-readable
//     ASCII table suitable for terminal output.
//   - Notifier: filters statuses that require attention and dispatches
//     a JSON alert payload to a configured webhook URL via HTTP POST.
package expiry
