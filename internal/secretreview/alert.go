package secretreview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// alertPayload is the JSON body sent to the webhook endpoint.
type alertPayload struct {
	GeneratedAt time.Time     `json:"generated_at"`
	AlertCount  int           `json:"alert_count"`
	Entries     []alertEntry  `json:"entries"`
}

// alertEntry represents a single overdue review in the payload.
type alertEntry struct {
	Path       string        `json:"path"`
	Reviewer   string        `json:"reviewer"`
	Interval   time.Duration `json:"interval_ns"`
	LastReview time.Time     `json:"last_review"`
	OverdueBy  time.Duration `json:"overdue_by_ns"`
	Status     string        `json:"status"`
}

// BuildAlertPayload constructs a webhook payload from the given review
// statuses, including only entries whose status is not OK.
func BuildAlertPayload(statuses []Status) alertPayload {
	var entries []alertEntry
	for _, s := range statuses {
		if s.OK {
			continue
		}
		entries = append(entries, alertEntry{
			Path:       s.Path,
			Reviewer:   s.Reviewer,
			Interval:   s.Interval,
			LastReview: s.LastReview,
			OverdueBy:  s.OverdueBy,
			Status:     statusLabel(s),
		})
	}
	return alertPayload{
		GeneratedAt: time.Now().UTC(),
		AlertCount:  len(entries),
		Entries:     entries,
	}
}

// SendAlert serialises the payload and POSTs it to the given webhook URL.
// It returns an error if the HTTP request fails or the server responds with
// a non-2xx status code.
func SendAlert(webhookURL string, payload alertPayload) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("secretreview: marshal payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("secretreview: send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("secretreview: webhook returned non-2xx status: %d", resp.StatusCode)
	}
	return nil
}
