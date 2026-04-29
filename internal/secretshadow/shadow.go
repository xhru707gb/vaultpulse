// Package secretshadow provides shadow-copy tracking for secrets,
// allowing detection of out-of-band changes by comparing live values
// against a stored shadow.
package secretshadow

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrPathEmpty is returned when an empty path is provided.
var ErrPathEmpty = errors.New("secretshadow: path must not be empty")

// ErrNotFound is returned when a path has no shadow entry.
var ErrNotFound = errors.New("secretshadow: path not found")

// Entry holds the shadow record for a single secret path.
type Entry struct {
	Path      string
	Hash      string
	CapturedAt time.Time
	Diverged  bool
}

// Tracker stores shadow hashes for secret paths.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	nowFn   func() time.Time
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		nowFn:   time.Now,
	}
}

func hashValue(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum)
}

// Capture stores a shadow hash for the given path and value.
func (t *Tracker) Capture(path, value string) error {
	if path == "" {
		return ErrPathEmpty
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[path] = Entry{
		Path:       path,
		Hash:       hashValue(value),
		CapturedAt: t.nowFn(),
		Diverged:   false,
	}
	return nil
}

// Check compares the live value against the stored shadow.
// Returns the Entry with Diverged set appropriately.
func (t *Tracker) Check(path, liveValue string) (Entry, error) {
	if path == "" {
		return Entry{}, ErrPathEmpty
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[path]
	if !ok {
		return Entry{}, ErrNotFound
	}
	e.Diverged = hashValue(liveValue) != e.Hash
	t.entries[path] = e
	return e, nil
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

// Remove deletes the shadow entry for the given path.
func (t *Tracker) Remove(path string) error {
	if path == "" {
		return ErrPathEmpty
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; !ok {
		return ErrNotFound
	}
	delete(t.entries, path)
	return nil
}
