// Package secretseal tracks whether secrets are sealed (write-protected)
// and provides evaluation of seal status per path.
package secretseal

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Status represents the seal evaluation result for a single secret path.
type Status struct {
	Path     string
	Sealed   bool
	SealedAt time.Time
	Reason   string
}

// entry holds internal state for a registered secret.
type entry struct {
	sealed   bool
	sealedAt time.Time
	reason   string
}

// Sealer manages seal state for secret paths.
type Sealer struct {
	mu      sync.RWMutex
	entries map[string]*entry
}

// New creates a new Sealer instance.
func New() *Sealer {
	return &Sealer{
		entries: make(map[string]*entry),
	}
}

// Seal marks a secret path as sealed with an optional reason.
func (s *Sealer) Seal(path, reason string, now time.Time) error {
	if path == "" {
		return errors.New("secretseal: path must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[path] = &entry{
		sealed:   true,
		sealedAt: now,
		reason:   reason,
	}
	return nil
}

// Unseal removes the seal from a secret path.
func (s *Sealer) Unseal(path string) error {
	if path == "" {
		return errors.New("secretseal: path must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[path]
	if !ok || !e.sealed {
		return fmt.Errorf("secretseal: path %q is not sealed", path)
	}
	e.sealed = false
	return nil
}

// Evaluate returns the seal Status for all registered paths.
func (s *Sealer) Evaluate() []Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Status, 0, len(s.entries))
	for path, e := range s.entries {
		out = append(out, Status{
			Path:     path,
			Sealed:   e.sealed,
			SealedAt: e.sealedAt,
			Reason:   e.reason,
		})
	}
	return out
}

// IsSealed reports whether the given path is currently sealed.
func (s *Sealer) IsSealed(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[path]
	return ok && e.sealed
}
