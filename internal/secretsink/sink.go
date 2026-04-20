// Package secretsink routes processed secret events to one or more
// registered output destinations (sinks), such as files, webhooks, or
// in-memory buffers.
package secretsink

import (
	"errors"
	"fmt"
	"sync"
)

// Event represents a secret lifecycle event delivered to sinks.
type Event struct {
	Path    string
	Kind    string // e.g. "expired", "rotated", "warning"
	Message string
}

// Sink is the interface that every output destination must satisfy.
type Sink interface {
	Name() string
	Send(Event) error
}

// Router dispatches events to all registered sinks.
type Router struct {
	mu    sync.RWMutex
	sinks map[string]Sink
}

// New creates an empty Router.
func New() *Router {
	return &Router{sinks: make(map[string]Sink)}
}

// Register adds a sink. Returns an error if the name is already taken.
func (r *Router) Register(s Sink) error {
	if s == nil {
		return errors.New("secretsink: sink must not be nil")
	}
	n := s.Name()
	if n == "" {
		return errors.New("secretsink: sink name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.sinks[n]; exists {
		return fmt.Errorf("secretsink: sink %q already registered", n)
	}
	r.sinks[n] = s
	return nil
}

// Deregister removes a sink by name. No-op if not found.
func (r *Router) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sinks, name)
}

// Dispatch sends the event to every registered sink, collecting errors.
func (r *Router) Dispatch(e Event) []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var errs []error
	for _, s := range r.sinks {
		if err := s.Send(e); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", s.Name(), err))
		}
	}
	return errs
}

// Len returns the number of registered sinks.
func (r *Router) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.sinks)
}
