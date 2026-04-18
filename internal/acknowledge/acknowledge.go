// Package acknowledge provides a mechanism to acknowledge alerts,
// suppressing repeated notifications for a known issue within a configurable window.
package acknowledge

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyAcknowledged is returned when a path is already acknowledged.
var ErrAlreadyAcknowledged = errors.New("acknowledge: path already acknowledged")

// ErrNotFound is returned when a path is not currently acknowledged.
var ErrNotFound = errors.New("acknowledge: path not found")

// Entry holds acknowledgement metadata for a secret path.
type Entry struct {
	Path       string
	AckedAt    time.Time
	ExpiresAt  time.Time
	AckedBy    string
}

// Tracker manages acknowledgements with expiry.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	now     func() time.Time
	window  time.Duration
}

// New creates a Tracker with the given acknowledgement window.
func New(window time.Duration) (*Tracker, error) {
	if window <= 0 {
		return nil, errors.New("acknowledge: window must be positive")
	}
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
		window:  window,
	}, nil
}

// Acknowledge records an acknowledgement for path by ackedBy.
func (t *Tracker) Acknowledge(path, ackedBy string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if e, ok := t.entries[path]; ok && e.ExpiresAt.After(now) {
		return ErrAlreadyAcknowledged
	}
	t.entries[path] = Entry{
		Path:      path,
		AckedAt:   now,
		ExpiresAt: now.Add(t.window),
		AckedBy:   ackedBy,
	}
	return nil
}

// IsAcknowledged reports whether path has an active acknowledgement.
func (t *Tracker) IsAcknowledged(path string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[path]
	return ok && e.ExpiresAt.After(t.now())
}

// Revoke removes an acknowledgement for path.
func (t *Tracker) Revoke(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; !ok {
		return ErrNotFound
	}
	delete(t.entries, path)
	return nil
}

// List returns all currently active acknowledgements.
func (t *Tracker) List() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		if e.ExpiresAt.After(now) {
			out = append(out, e)
		}
	}
	return out
}
