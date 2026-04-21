// Package secretarchive provides archival tracking for secrets that have
// been rotated out or expired, preserving a historical record.
package secretarchive

import (
	"errors"
	"sync"
	"time"
)

// Entry represents an archived secret record.
type Entry struct {
	Path       string
	Version    int
	ArchivedAt time.Time
	Reason     string
}

// Archiver stores archived secret entries.
type Archiver struct {
	mu      sync.RWMutex
	entries []Entry
	now     func() time.Time
}

// New creates a new Archiver.
func New() *Archiver {
	return &Archiver{now: time.Now}
}

// Archive records a secret path and version as archived with a given reason.
func (a *Archiver) Archive(path string, version int, reason string) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if version < 1 {
		return errors.New("version must be >= 1")
	}
	if reason == "" {
		return errors.New("reason must not be empty")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, Entry{
		Path:       path,
		Version:    version,
		ArchivedAt: a.now().UTC(),
		Reason:     reason,
	})
	return nil
}

// All returns a copy of all archived entries.
func (a *Archiver) All() []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Entry, len(a.entries))
	copy(out, a.entries)
	return out
}

// ForPath returns all archived entries for a given path.
func (a *Archiver) ForPath(path string) []Entry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var out []Entry
	for _, e := range a.entries {
		if e.Path == path {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the total number of archived entries.
func (a *Archiver) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}
