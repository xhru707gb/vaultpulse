// Package snapshot provides functionality to capture and persist the current
// state of Vault secret expiry and rotation statuses to a JSON file.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yourusername/vaultpulse/internal/expiry"
	"github.com/yourusername/vaultpulse/internal/rotation"
)

// Snapshot holds a point-in-time capture of secret health.
type Snapshot struct {
	CapturedAt     time.Time               `json:"captured_at"`
	ExpiryStatuses []expiry.Status         `json:"expiry_statuses"`
	RotationStatuses []rotation.Status     `json:"rotation_statuses"`
}

// Writer writes snapshots to a file path.
type Writer struct {
	path string
	now  func() time.Time
}

// NewWriter creates a Writer that persists snapshots to the given file path.
func NewWriter(path string) *Writer {
	return &Writer{path: path, now: time.Now}
}

// Write serialises the given statuses into a Snapshot and writes it to disk.
func (w *Writer) Write(es []expiry.Status, rs []rotation.Status) error {
	snap := Snapshot{
		CapturedAt:       w.now().UTC(),
		ExpiryStatuses:   es,
		RotationStatuses: rs,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(w.path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}

	return nil
}

// Load reads and deserialises a Snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}

	return &snap, nil
}
