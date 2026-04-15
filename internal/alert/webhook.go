package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/vaultpulse/internal/expiry"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp string          `json:"timestamp"`
	Alerts    []AlertEntry    `json:"alerts"`
}

// AlertEntry represents a single secret alert in the payload.
type AlertEntry struct {
	Path   string `json:"path"`
	Status string `json:"status"`
	TTL    string `json:"ttl"`
}

// buildPayload constructs a WebhookPayload from the given statuses,
// including only entries that are in Warning or Expired state.
func buildPayload(statuses []expiry.SecretStatus) WebhookPayload {
	entries := make([]AlertEntry, 0, len(statuses))
	for _, s := range statuses {
		if s.State == expiry.StateOK {
			continue
		}
		entries = append(entries, AlertEntry{
			Path:   s.Path,
			Status: string(s.State),
			TTL:    formatTTLDuration(s.TTL),
		})
	}
	return WebhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Alerts:    entries,
	}
}

// formatTTLDuration formats a duration in seconds as a human-readable string.
func formatTTLDuration(seconds int64) string {
	if seconds <= 0 {
		return "expired"
	}
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// sendWebhook marshals the payload and POSTs it to the given URL.
func sendWebhook(url string, payload WebhookPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal webhook payload: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("post webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
