package secretexpiry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertPayload is the JSON body sent to the webhook for expiry alerts.
type AlertPayload struct {
	GeneratedAt time.Time     `json:"generated_at"`
	AlertCount  int           `json:"alert_count"`
	Entries     []alertEntry  `json:"entries"`
}

type alertEntry struct {
	Path      string `json:"path"`
	State     string `json:"state"`
	Remaining string `json:"remaining"`
}

// BuildAlertPayload returns a payload containing only warning or expired entries.
func BuildAlertPayload(statuses []Status) AlertPayload {
	var entries []alertEntry
	for _, s := range statuses {
		if s.State == StateOK {
			continue
		}
		entries = append(entries, alertEntry{
			Path:      s.Path,
			State:     stateLabel(s.State),
			Remaining: formatRemaining(s.Remaining),
		})
	}
	return AlertPayload{
		GeneratedAt: time.Now().UTC(),
		AlertCount:  len(entries),
		Entries:     entries,
	}
}

// SendAlert marshals the payload and posts it to webhookURL.
func SendAlert(webhookURL string, payload AlertPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("secretexpiry: marshal payload: %w", err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("secretexpiry: send alert: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("secretexpiry: unexpected status %d", resp.StatusCode)
	}
	return nil
}
