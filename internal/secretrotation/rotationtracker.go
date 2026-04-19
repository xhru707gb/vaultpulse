// Package secretrotation tracks per-secret rotation history and due dates.
package secretrotation

import (
	"errors"
	"sync"
	"time"
)

// Entry holds rotation metadata for a single secret path.
type Entry struct {
	Path        string
	Interval    time.Duration
	LastRotated time.Time
	NextDue     time.Time
	Overdue     bool
}

// Tracker manages rotation records for multiple secrets.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a Tracker with an optional clock override.
func New(now func() time.Time) (*Tracker, error) {
	if now == nil {
		now = time.Now
	}
	return &Tracker{entries: make(map[string]Entry), now: now}, nil
}

// Register adds or updates a secret rotation record.
func (t *Tracker) Register(path string, interval time.Duration, lastRotated time.Time) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if interval <= 0 {
		return errors.New("interval must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	nextDue := lastRotated.Add(interval)
	t.entries[path] = Entry{
		Path:        path,
		Interval:    interval,
		LastRotated: lastRotated,
		NextDue:     nextDue,
		Overdue:     t.now().After(nextDue),
	}
	return nil
}

// Get returns the rotation entry for a path.
func (t *Tracker) Get(path string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	return e, ok
}

// All returns a snapshot of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// OverdueCount returns the number of overdue secrets.
func (t *Tracker) OverdueCount() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	count := 0
	for _, e := range t.entries {
		if e.Overdue {
			count++
		}
	}
	return count
}
