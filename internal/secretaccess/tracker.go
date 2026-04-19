// Package secretaccess tracks how frequently secrets are accessed
// and surfaces paths that may be over- or under-utilised.
package secretaccess

import (
	"errors"
	"sync"
	"time"
)

// Entry records access statistics for a single secret path.
type Entry struct {
	Path        string
	AccessCount int
	LastAccess  time.Time
}

// Tracker maintains access counts for secret paths.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// Record increments the access counter for path.
func (t *Tracker) Record(path string) error {
	if path == "" {
		return errors.New("secretaccess: path must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[path]
	if !ok {
		e = &Entry{Path: path}
		t.entries[path] = e
	}
	e.AccessCount++
	e.LastAccess = t.now()
	return nil
}

// Get returns the Entry for path and whether it exists.
func (t *Tracker) Get(path string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[path]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all recorded access data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*Entry)
}
