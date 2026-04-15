package rotation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AlertPayload is the JSON body sent to a webhook for overdue rotations.
type AlertPayload struct {
	Timestamp string          `json:"timestamp"`
	Overdue   []OverdueEntry  `json:"overdue"`
}

// OverdueEntry represents a single secret path that is overdue for rotation.
type OverdueEntry struct {
	Path     string `json:"path"`
	Interval string `json:"interval"`
	Overdue  string `json:"overdue_by"`
}

// BuildAlertPayload constructs an AlertPayload from a slice of EvaluationStatus,
// including only paths whose Overdue flag is true.
func BuildAlertPayload(statuses []EvaluationStatus) AlertPayload {
	p := AlertPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Overdue:   []OverdueEntry{},
	}
	for _, s := range statuses {
		if !s.Overdue {
			continue
		}
		p.Overdue = append(p.Overdue, OverdueEntry{
			Path:     s.Path,
			Interval: s.Interval.String(),
			Overdue:  formatDueIn(-s.DueIn),
		})
	}
	return p
}

// SendAlert posts the AlertPayload to the given webhookURL.
// It returns an error if the payload cannot be marshalled or if the server
// responds with a non-2xx status code.
func SendAlert(webhookURL string, payload AlertPayload) error {
	if len(payload.Overdue) == 0 {
		return nil
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("rotation alert: marshal payload: %w", err)
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("rotation alert: post webhook: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("rotation alert: unexpected status %d", resp.StatusCode)
	}
	return nil
}
