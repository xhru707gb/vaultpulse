// Package secretlease tracks Vault dynamic secret leases and warns when
// they are approaching expiry or have already expired.
package secretlease

import (
	"errors"
	"sync"
	"time"
)

// State represents the lease health of a secret.
type State int

const (
	StateOK State = iota
	StateWarning
	StateExpired
)

// Entry holds lease metadata for a single secret path.
type Entry struct {
	Path      string
	LeaseID   string
	ExpiresAt time.Time
	Renewable bool
}

// Status is the evaluated result for a lease.
type Status struct {
	Entry
	State     State
	Remaining time.Duration
}

// Watcher monitors registered leases and evaluates their expiry state.
type Watcher struct {
	mu      sync.RWMutex
	entries map[string]Entry
	warn    time.Duration
	now     func() time.Time
}

// New creates a Watcher with the given warning threshold.
func New(warnBefore time.Duration) (*Watcher, error) {
	if warnBefore <= 0 {
		return nil, errors.New("secretlease: warnBefore must be positive")
	}
	return &Watcher{
		entries: make(map[string]Entry),
		warn:    warnBefore,
		now:     time.Now,
	}, nil
}

// Register adds or replaces a lease entry for the given path.
func (w *Watcher) Register(e Entry) error {
	if e.Path == "" {
		return errors.New("secretlease: path must not be empty")
	}
	if e.LeaseID == "" {
		return errors.New("secretlease: leaseID must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[e.Path] = e
	return nil
}

// Remove deletes the lease entry for the given path.
func (w *Watcher) Remove(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, path)
}

// Evaluate returns the current Status for every registered lease.
func (w *Watcher) Evaluate() []Status {
	w.mu.RLock()
	defer w.mu.RUnlock()
	now := w.now()
	out := make([]Status, 0, len(w.entries))
	for _, e := range w.entries {
		remaining := e.ExpiresAt.Sub(now)
		var state State
		switch {
		case remaining <= 0:
			state = StateExpired
		case remaining <= w.warn:
			state = StateWarning
		default:
			state = StateOK
		}
		out = append(out, Status{Entry: e, State: state, Remaining: remaining})
	}
	return out
}
