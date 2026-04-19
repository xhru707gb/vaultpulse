// Package secretaudit tracks access and mutation events for secrets,
// providing a tamper-evident trail of who touched what and when.
package secretaudit

import (
	"errors"
	"sync"
	"time"
)

// EventKind describes the type of secret event.
type EventKind string

const (
	EventRead   EventKind = "read"
	EventWrite  EventKind = "write"
	EventDelete EventKind = "delete"
	EventRotate EventKind = "rotate"
)

// Event represents a single audit event for a secret path.
type Event struct {
	Path      string
	Kind      EventKind
	Actor     string
	Timestamp time.Time
}

// Auditor records and retrieves secret audit events.
type Auditor struct {
	mu     sync.RWMutex
	events []Event
	now    func() time.Time
}

// New creates a new Auditor.
func New() *Auditor {
	return &Auditor{now: time.Now}
}

// Record appends an audit event.
func (a *Auditor) Record(path string, kind EventKind, actor string) error {
	if path == "" {
		return errors.New("secretaudit: path must not be empty")
	}
	if actor == "" {
		return errors.New("secretaudit: actor must not be empty")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, Event{
		Path:      path,
		Kind:      kind,
		Actor:     actor,
		Timestamp: a.now(),
	})
	return nil
}

// ForPath returns all events recorded for the given path.
func (a *Auditor) ForPath(path string) []Event {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var out []Event
	for _, e := range a.events {
		if e.Path == path {
			out = append(out, e)
		}
	}
	return out
}

// All returns a copy of all recorded events.
func (a *Auditor) All() []Event {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Event, len(a.events))
	copy(out, a.events)
	return out
}

// Reset clears all recorded events.
func (a *Auditor) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = nil
}
