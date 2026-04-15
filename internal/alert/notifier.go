// Package alert provides alerting hooks for secret expiry notifications.
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/vaultpulse/internal/expiry"
)

// WebhookPayload is the JSON body sent to a webhook endpoint.
type WebhookPayload struct {
	Timestamp time.Time        `json:"timestamp"`
	Alerts    []expiry.Status  `json:"alerts"`
	Summary   string           `json:"summary"`
}

// Notifier sends alerts for expiring or expired secrets.
type Notifier struct {
	webhookURL string
	httpClient *http.Client
}

// NewNotifier creates a Notifier with the given webhook URL.
func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify filters statuses that require alerting and posts them to the webhook.
// Only Warning and Expired states are dispatched.
func (n *Notifier) Notify(statuses []expiry.Status) error {
	var alerts []expiry.Status
	for _, s := range statuses {
		if s.State == expiry.StateWarning || s.State == expiry.StateExpired {
			alerts = append(alerts, s)
		}
	}
	if len(alerts) == 0 {
		return nil
	}

	payload := WebhookPayload{
		Timestamp: time.Now().UTC(),
		Alerts:    alerts,
		Summary:   fmt.Sprintf("%d secret(s) require attention", len(alerts)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("alert: marshal payload: %w", err)
	}

	resp, err := n.httpClient.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("alert: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alert: webhook returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
