// Package secretage tracks how long secrets have been in use
// and flags those that have exceeded a configurable maximum age.
package secretage

import (
	"errors"
	"sync"
	"time"
)

// Entry holds age metadata for a single secret path.
type Entry struct {
	Path      string
	CreatedAt time.Time
	MaxAge    time.Duration
}

// Status is the evaluated age result for a secret.
type Status struct {
	Path    string
	Age     time.Duration
	MaxAge  time.Duration
	Overdue bool
}

// Tracker manages secret age entries.
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

// Register adds or updates a secret entry.
func (t *Tracker) Register(path string, createdAt time.Time, maxAge time.Duration) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if maxAge <= 0 {
		return errors.New("maxAge must be positive")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[path] = Entry{Path: path, CreatedAt: createdAt, MaxAge: maxAge}
	return nil
}

// Evaluate returns the age status for a single path.
func (t *Tracker) Evaluate(path string) (Status, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	if !ok {
		return Status{}, false
	}
	age := t.now().Sub(e.CreatedAt)
	return Status{Path: path, Age: age, MaxAge: e.MaxAge, Overdue: age > e.MaxAge}, true
}

// EvaluateAll returns statuses for all registered secrets.
func (t *Tracker) EvaluateAll() []Status {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Status, 0, len(t.entries))
	for _, e := range t.entries {
		age := t.now().Sub(e.CreatedAt)
		out = append(out, Status{Path: e.Path, Age: age, MaxAge: e.MaxAge, Overdue: age > e.MaxAge})
	}
	return out
}
