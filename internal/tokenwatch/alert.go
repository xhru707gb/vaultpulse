package tokenwatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertPayload is the webhook body sent when token states require attention.
type AlertPayload struct {
	Timestamp time.Time     `json:"timestamp"`
	Alerts    []TokenAlert  `json:"alerts"`
}

// TokenAlert represents a single token entry that is not OK.
type TokenAlert struct {
	Path      string `json:"path"`
	State     string `json:"state"`
	Remaining string `json:"remaining"`
}

// BuildAlertPayload filters statuses to those that are not StateOK and
// returns a payload ready for dispatch.
func BuildAlertPayload(statuses []Status, now time.Time) AlertPayload {
	var alerts []TokenAlert
	for _, s := range statuses {
		if s.State == StateOK {
			continue
		}
		alerts = append(alerts, TokenAlert{
			Path:      s.Path,
			State:     stateLabel(s.State),
			Remaining: formatRemaining(s.Remaining),
		})
	}
	return AlertPayload{Timestamp: now.UTC(), Alerts: alerts}
}

// SendAlert posts the payload to webhookURL. It returns an error if the
// request fails or the server responds with a non-2xx status.
func SendAlert(webhookURL string, payload AlertPayload) error {
	if len(payload.Alerts) == 0 {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("tokenwatch: marshal payload: %w", err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("tokenwatch: send alert: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("tokenwatch: unexpected status %d", resp.StatusCode)
	}
	return nil
}
