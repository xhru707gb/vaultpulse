// Package secretexpiry tracks and evaluates secret expiry windows.
package secretexpiry

import (
	"errors"
	"sync"
	"time"
)

// Entry holds expiry metadata for a single secret path.
type Entry struct {
	Path      string
	ExpiresAt time.Time
	WarnBefore time.Duration
}

// Status is the evaluated result for a secret.
type Status struct {
	Entry
	Expired  bool
	Warning  bool
	Remaining time.Duration
}

// Tracker manages secret expiry entries.
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

// Register adds or updates a secret expiry entry.
func (t *Tracker) Register(e Entry) error {
	if e.Path == "" {
		return errors.New("path must not be empty")
	}
	if e.ExpiresAt.IsZero() {
		return errors.New("expiresAt must not be zero")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[e.Path] = e
	return nil
}

// Evaluate returns the Status for a given path.
func (t *Tracker) Evaluate(path string) (Status, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	if !ok {
		return Status{}, false
	}
	now := t.now()
	remaining := e.ExpiresAt.Sub(now)
	return Status{
		Entry:     e,
		Expired:   now.After(e.ExpiresAt),
		Warning:   !now.After(e.ExpiresAt) && remaining <= e.WarnBefore,
		Remaining: remaining,
	}, true
}

// EvaluateAll returns statuses for all registered secrets.
func (t *Tracker) EvaluateAll() []Status {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Status, 0, len(t.entries))
	now := t.now()
	for _, e := range t.entries {
		remaining := e.ExpiresAt.Sub(now)
		out = append(out, Status{
			Entry:     e,
			Expired:   now.After(e.ExpiresAt),
			Warning:   !now.After(e.ExpiresAt) && remaining <= e.WarnBefore,
			Remaining: remaining,
		})
	}
	return out
}
