// Package cooldown provides a simple per-key cooldown tracker that prevents
// repeated actions within a configurable quiet period.
package cooldown

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidWindow is returned when the cooldown window is not positive.
var ErrInvalidWindow = errors.New("cooldown: window must be greater than zero")

// entry holds the expiry time for a single key.
type entry struct {
	expiresAt time.Time
}

// Tracker tracks per-key cooldown periods.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	now    func() time.Time
	keys   map[string]entry
}

// New creates a Tracker with the given cooldown window.
func New(window time.Duration) (*Tracker, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Tracker{
		window: window,
		now:    time.Now,
		keys:   make(map[string]entry),
	}, nil
}

// IsCoolingDown returns true if the key is still within its cooldown period.
func (t *Tracker) IsCoolingDown(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.keys[key]
	if !ok {
		return false
	}
	return t.now().Before(e.expiresAt)
}

// Record marks the key as active, resetting its cooldown window.
func (t *Tracker) Record(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.keys[key] = entry{expiresAt: t.now().Add(t.window)}
}

// Reset removes the cooldown entry for the given key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.keys, key)
}

// Len returns the number of tracked keys.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.keys)
}
