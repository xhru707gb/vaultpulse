// Package quarantine tracks secrets flagged for immediate rotation or revocation.
package quarantine

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyQuarantined is returned when a secret is already in quarantine.
var ErrAlreadyQuarantined = errors.New("secret is already quarantined")

// ErrNotFound is returned when a secret is not in quarantine.
var ErrNotFound = errors.New("secret not found in quarantine")

// Reason describes why a secret was quarantined.
type Reason string

const (
	ReasonExpired  Reason = "expired"
	ReasonLeaked   Reason = "leaked"
	ReasonPolicy   Reason = "policy_violation"
	ReasonManual   Reason = "manual"
)

// Entry holds quarantine metadata for a single secret path.
type Entry struct {
	Path        string
	Reason      Reason
	QuarantinedAt time.Time
	Note        string
}

// Store holds quarantined secret paths.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns an initialised Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Add adds a secret path to quarantine. Returns ErrAlreadyQuarantined if present.
func (s *Store) Add(path string, reason Reason, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[path]; ok {
		return ErrAlreadyQuarantined
	}
	s.entries[path] = Entry{
		Path:          path,
		Reason:        reason,
		QuarantinedAt: s.now(),
		Note:          note,
	}
	return nil
}

// Remove removes a secret path from quarantine.
func (s *Store) Remove(path string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[path]; !ok {
		return ErrNotFound
	}
	delete(s.entries, path)
	return nil
}

// IsQuarantined reports whether a path is currently quarantined.
func (s *Store) IsQuarantined(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.entries[path]
	return ok
}

// All returns a snapshot of all quarantined entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
