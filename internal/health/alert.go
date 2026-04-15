package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertPayload is the JSON body sent to a webhook for unhealthy Vault nodes.
type AlertPayload struct {
	Timestamp string        `json:"timestamp"`
	Alerts    []AlertEntry  `json:"alerts"`
}

// AlertEntry represents a single unhealthy node in the payload.
type AlertEntry struct {
	Node      string `json:"node"`
	Status    string `json:"status"`
	Sealed    bool   `json:"sealed"`
	LatencyMs int64  `json:"latency_ms"`
}

// BuildAlertPayload filters statuses to only unhealthy or sealed nodes and
// returns a payload ready for webhook delivery. Returns nil if no alerts.
func BuildAlertPayload(statuses []Status) *AlertPayload {
	var entries []AlertEntry
	for _, s := range statuses {
		if s.Healthy && !s.Sealed {
			continue
		}
		entries = append(entries, AlertEntry{
			Node:      s.Node,
			Status:    statusLabel(s),
			Sealed:    s.Sealed,
			LatencyMs: s.LatencyMs,
		})
	}
	if len(entries) == 0 {
		return nil
	}
	return &AlertPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Alerts:    entries,
	}
}

// SendAlert posts the payload to webhookURL. Returns an error if the
// server responds with a non-2xx status code or the request fails.
func SendAlert(webhookURL string, payload *AlertPayload) error {
	if payload == nil {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("health alert: marshal payload: %w", err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("health alert: send webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("health alert: webhook returned %d", resp.StatusCode)
	}
	return nil
}
