// Package grace tracks secrets within a configurable grace period after expiry.
package grace

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyTracked is returned when a path is registered more than once.
var ErrAlreadyTracked = errors.New("grace: path already tracked")

// ErrNotFound is returned when a path is not registered.
var ErrNotFound = errors.New("grace: path not found")

// Entry holds the expiry and grace deadline for a secret path.
type Entry struct {
	Path         string
	ExpiredAt    time.Time
	GraceEndsAt  time.Time
}

// InGrace reports whether the entry is within the grace window.
func (e Entry) InGrace(now time.Time) bool {
	return now.After(e.ExpiredAt) && !now.After(e.GraceEndsAt)
}

// Tracker manages grace-period entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	window  time.Duration
	now     func() time.Time
}

// New creates a Tracker with the given grace window.
func New(window time.Duration, nowFn func() time.Time) (*Tracker, error) {
	if window <= 0 {
		return nil, errors.New("grace: window must be positive")
	}
	if nowFn == nil {
		nowFn = time.Now
	}
	return &Tracker{entries: make(map[string]Entry), window: window, now: nowFn}, nil
}

// Register records a secret expiry and computes the grace deadline.
func (t *Tracker) Register(path string, expiredAt time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; ok {
		return ErrAlreadyTracked
	}
	t.entries[path] = Entry{
		Path:        path,
		ExpiredAt:   expiredAt,
		GraceEndsAt: expiredAt.Add(t.window),
	}
	return nil
}

// Remove deletes a tracked entry.
func (t *Tracker) Remove(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; !ok {
		return ErrNotFound
	}
	delete(t.entries, path)
	return nil
}

// Active returns all entries currently within the grace period.
func (t *Tracker) Active() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if e.InGrace(now) {
			out = append(out, e)
		}
	}
	return out
}
