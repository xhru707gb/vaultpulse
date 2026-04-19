// Package secretlifecycle tracks the full lifecycle state of a secret,
// combining expiry, rotation, and access metadata into a unified status.
package secretlifecycle

import (
	"errors"
	"sync"
	"time"
)

// State represents the overall lifecycle state of a secret.
type State string

const (
	StateActive   State = "active"
	StateExpiring State = "expiring"
	StateExpired  State = "expired"
	StateStale    State = "stale"
)

// Entry holds lifecycle metadata for a single secret path.
type Entry struct {
	Path        string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	LastRotated time.Time
	LastAccess  time.Time
	MaxAge      time.Duration
	WarnBefore  time.Duration
}

// Status is the evaluated lifecycle status of a secret.
type Status struct {
	Entry
	State   State
	Message string
}

// Tracker manages lifecycle entries and evaluates their states.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Register adds or replaces a lifecycle entry.
func (t *Tracker) Register(e Entry) error {
	if e.Path == "" {
		return errors.New("path must not be empty")
	}
	if e.MaxAge <= 0 {
		return errors.New("max age must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[e.Path] = e
	return nil
}

// Evaluate returns the current Status for the given path.
func (t *Tracker) Evaluate(path string) (Status, bool) {
	t.mu.RLock()
	e, ok := t.entries[path]
	t.mu.RUnlock()
	if !ok {
		return Status{}, false
	}
	now := t.now()
	s := Status{Entry: e}
	switch {
	case !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt):
		s.State = StateExpired
		s.Message = "secret has expired"
	case !e.ExpiresAt.IsZero() && now.After(e.ExpiresAt.Add(-e.WarnBefore)):
		s.State = StateExpiring
		s.Message = "secret is expiring soon"
	case !e.LastRotated.IsZero() && now.Sub(e.LastRotated) > e.MaxAge:
		s.State = StateStale
		s.Message = "secret has not been rotated within max age"
	default:
		s.State = StateActive
		s.Message = "secret is active"
	}
	return s, true
}

// EvaluateAll returns statuses for all registered entries.
func (t *Tracker) EvaluateAll() []Status {
	t.mu.RLock()
	paths := make([]string, 0, len(t.entries))
	for p := range t.entries {
		paths = append(paths, p)
	}
	t.mu.RUnlock()
	out := make([]Status, 0, len(paths))
	for _, p := range paths {
		if s, ok := t.Evaluate(p); ok {
			out = append(out, s)
		}
	}
	return out
}
