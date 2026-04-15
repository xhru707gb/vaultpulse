// Package alert implements alerting hooks for the vaultpulse CLI.
//
// It provides a Notifier that inspects a slice of expiry.Status values and
// dispatches a JSON webhook payload for any secrets in Warning or Expired
// state. Healthy secrets (StateOK) are silently ignored so that only
// actionable events reach downstream consumers such as Slack, PagerDuty,
// or a custom HTTP endpoint.
//
// Basic usage:
//
//	n := alert.NewNotifier("https://hooks.example.com/vault-alerts")
//	if err := n.Notify(statuses); err != nil {
//		log.Printf("alert delivery failed: %v", err)
//	}
package alert
