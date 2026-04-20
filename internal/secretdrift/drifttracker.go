// Package secretdrift detects and tracks value drift between secret snapshots.
package secretdrift

import (
	"errors"
	"sync"
	"time"
)

// DriftEntry records a detected drift for a secret path.
type DriftEntry struct {
	Path        string
	PreviousHash string
	CurrentHash  string
	DetectedAt  time.Time
}

// Tracker holds drift state for monitored secrets.
type Tracker struct {
	mu      sync.RWMutex
	hashes  map[string]string
	drifts  []DriftEntry
	now     func() time.Time
}

// New creates a Tracker with an optional clock override.
func New(nowFn func() time.Time) (*Tracker, error) {
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Tracker{
		hashes: make(map[string]string),
		now:    nowFn,
	}, nil
}

// Record compares the incoming hash against the stored baseline.
// If a drift is detected it is appended to the drift log.
func (t *Tracker) Record(path, hash string) error {
	if path == "" {
		return errors.New("secretdrift: path must not be empty")
	}
	if hash == "" {
		return errors.New("secretdrift: hash must not be empty")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	prev, exists := t.hashes[path]
	t.hashes[path] = hash

	if exists && prev != hash {
		t.drifts = append(t.drifts, DriftEntry{
			Path:         path,
			PreviousHash: prev,
			CurrentHash:  hash,
			DetectedAt:   t.now(),
		})
	}
	return nil
}

// Drifts returns all recorded drift entries.
func (t *Tracker) Drifts() []DriftEntry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]DriftEntry, len(t.drifts))
	copy(out, t.drifts)
	return out
}

// Reset clears the drift log but retains baseline hashes.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.drifts = nil
}
