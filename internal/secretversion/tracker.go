package secretversion

import (
	"errors"
	"sync"
	"time"
)

// Entry holds version metadata for a secret path.
type Entry struct {
	Path      string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Tracker maintains a registry of secret versions.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Register adds or updates a secret version entry.
func (t *Tracker) Register(path string, version int) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if version < 1 {
		return errors.New("version must be >= 1")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if existing, ok := t.entries[path]; ok {
		existing.Version = version
		existing.UpdatedAt = now
		t.entries[path] = existing
	} else {
		t.entries[path] = Entry{
			Path:      path,
			Version:   version,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return nil
}

// Get returns the entry for a path.
func (t *Tracker) Get(path string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	return e, ok
}

// All returns a snapshot of all entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// Remove deletes an entry by path.
func (t *Tracker) Remove(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; !ok {
		return errors.New("path not found")
	}
	delete(t.entries, path)
	return nil
}
