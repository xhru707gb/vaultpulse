// Package suppress provides a mechanism to suppress repeated alerts
// for a given path within a configurable time window.
package suppress

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidWindow is returned when the suppression window is non-positive.
var ErrInvalidWindow = errors.New("suppress: window must be greater than zero")

// Record holds suppression state for a single key.
type Record struct {
	SuppressedUntil time.Time
	Count           int
}

// Suppressor tracks suppressed alert keys.
type Suppressor struct {
	mu     sync.Mutex
	window time.Duration
	now    func() time.Time
	state  map[string]Record
}

// New creates a Suppressor with the given suppression window.
func New(window time.Duration) (*Suppressor, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Suppressor{
		window: window,
		now:    time.Now,
		state:  make(map[string]Record),
	}, nil
}

// IsSuppressed returns true if the key is currently suppressed.
func (s *Suppressor) IsSuppressed(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.state[key]
	if !ok {
		return false
	}
	return s.now().Before(r.SuppressedUntil)
}

// Record marks the key as suppressed for the configured window.
func (s *Suppressor) Record(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.state[key]
	r.SuppressedUntil = s.now().Add(s.window)
	r.Count++
	s.state[key] = r
}

// Reset clears suppression state for a key.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.state, key)
}

// All returns a snapshot of all suppression records.
func (s *Suppressor) All() map[string]Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]Record, len(s.state))
	for k, v := range s.state {
		out[k] = v
	}
	return out
}
