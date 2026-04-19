// Package ownership tracks secret path ownership and contact metadata.
package ownership

import (
	"errors"
	"sync"
	"time"
)

// Entry holds ownership metadata for a secret path.
type Entry struct {
	Path      string
	Owner     string
	Team      string
	Contact   string
	CreatedAt time.Time
}

// Registry stores ownership entries keyed by secret path.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a new Registry.
func New() *Registry {
	return &Registry{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Register adds or replaces ownership for the given path.
func (r *Registry) Register(path, owner, team, contact string) error {
	if path == "" {
		return errors.New("ownership: path must not be empty")
	}
	if owner == "" {
		return errors.New("ownership: owner must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[path] = Entry{
		Path:      path,
		Owner:     owner,
		Team:      team,
		Contact:   contact,
		CreatedAt: r.now(),
	}
	return nil
}

// Get returns the ownership entry for path, or false if not found.
func (r *Registry) Get(path string) (Entry, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[path]
	return e, ok
}

// Remove deletes the ownership record for path.
func (r *Registry) Remove(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[path]; !ok {
		return errors.New("ownership: path not found")
	}
	delete(r.entries, path)
	return nil
}

// All returns a snapshot of all entries.
func (r *Registry) All() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}
