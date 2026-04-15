// Package health provides Vault node health checking, formatting,
// and alerting capabilities for the vaultpulse CLI.
//
// It exposes:
//   - Checker: queries the Vault /sys/health endpoint and returns
//     a slice of Status values describing each node's state.
//   - FormatTable: renders health statuses as a human-readable
//     ASCII table suitable for terminal output.
//   - BuildAlertPayload / SendAlert: filter unhealthy nodes and
//     POST a JSON alert payload to a configured webhook URL.
package health
