// Package secretmap provides a registry for tracking secret metadata
// indexed by path, supporting lookup, listing, and removal.
package secretmap

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when a secret path is not registered.
var ErrNotFound = errors.New("secretmap: path not found")

// ErrAlreadyExists is returned when a path is registered twice.
var ErrAlreadyExists = errors.New("secretmap: path already exists")

// Entry holds metadata for a single secret.
type Entry struct {
	Path        string
	Version     int
	ExpiresAt   time.Time
	LastRotated time.Time
	Owner       string
}

// Registry stores secret entries indexed by path.
type Registry struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{entries: make(map[string]Entry)}
}

// Register adds an entry. Returns ErrAlreadyExists if the path is taken.
func (r *Registry) Register(e Entry) error {
	if e.Path == "" {
		return errors.New("secretmap: path must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[e.Path]; ok {
		return ErrAlreadyExists
	}
	r.entries[e.Path] = e
	return nil
}

// Get returns the entry for path. Returns ErrNotFound if absent.
func (r *Registry) Get(path string) (Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.entries[path]
	if !ok {
		return Entry{}, ErrNotFound
	}
	return e, nil
}

// Remove deletes the entry for path. Returns ErrNotFound if absent.
func (r *Registry) Remove(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[path]; !ok {
		return ErrNotFound
	}
	delete(r.entries, path)
	return nil
}

// All returns a snapshot of all registered entries.
func (r *Registry) All() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Entry, 0, len(r.entries))
	for _, e := range r.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of registered entries.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}
