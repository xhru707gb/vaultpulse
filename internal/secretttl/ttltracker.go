// Package secretttl tracks remaining TTL for secrets and classifies
// each entry as OK, Warning, or Expired based on configurable thresholds.
package secretttl

import (
	"errors"
	"sync"
	"time"
)

// State represents the TTL health of a secret.
type State int

const (
	StateOK      State = iota
	StateWarning       // TTL is below the warning threshold
	StateExpired       // TTL has elapsed
)

// Entry holds the registration data for a single secret.
type Entry struct {
	Path      string
	ExpiresAt time.Time
	WarningIn time.Duration // warn when remaining TTL drops below this
}

// Status is the evaluated result for a secret.
type Status struct {
	Path      string
	Remaining time.Duration
	State     State
}

// Tracker manages TTL entries and evaluates their current state.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a Tracker with an optional clock override (nil = time.Now).
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		entries: make(map[string]Entry),
		now:     now,
	}
}

// Register adds or replaces a TTL entry for the given secret path.
func (t *Tracker) Register(e Entry) error {
	if e.Path == "" {
		return errors.New("secretttl: path must not be empty")
	}
	if e.ExpiresAt.IsZero() {
		return errors.New("secretttl: expiresAt must not be zero")
	}
	if e.WarningIn <= 0 {
		return errors.New("secretttl: warningIn must be positive")
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
	return t.evaluate(e), true
}

// EvaluateAll returns statuses for every registered secret.
func (t *Tracker) EvaluateAll() []Status {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Status, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, t.evaluate(e))
	}
	return out
}

func (t *Tracker) evaluate(e Entry) Status {
	now := t.now()
	remaining := e.ExpiresAt.Sub(now)
	var state State
	switch {
	case remaining <= 0:
		state = StateExpired
	case remaining <= e.WarningIn:
		state = StateWarning
	default:
		state = StateOK
	}
	return Status{Path: e.Path, Remaining: remaining, State: state}
}
