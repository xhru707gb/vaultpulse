// Package secretrelay provides a fan-out relay that forwards secret events
// to multiple downstream consumers via registered handler functions.
package secretrelay

import (
	"errors"
	"fmt"
	"sync"
)

// Handler is a function that receives a secret path and its associated payload.
type Handler func(path string, payload map[string]string) error

// Relay fans out secret events to all registered handlers.
type Relay struct {
	mu       sync.RWMutex
	handlers map[string]Handler
}

// New creates a new Relay.
func New() *Relay {
	return &Relay{
		handlers: make(map[string]Handler),
	}
}

// Register adds a named handler to the relay.
// Returns an error if the name is empty, the handler is nil, or the name is
// already registered.
func (r *Relay) Register(name string, h Handler) error {
	if name == "" {
		return errors.New("secretrelay: handler name must not be empty")
	}
	if h == nil {
		return errors.New("secretrelay: handler must not be nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.handlers[name]; exists {
		return fmt.Errorf("secretrelay: handler %q already registered", name)
	}
	r.handlers[name] = h
	return nil
}

// Deregister removes a handler by name. No-op if the name is unknown.
func (r *Relay) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.handlers, name)
}

// Dispatch sends the event to all registered handlers and collects any errors.
// All handlers are invoked even if one returns an error.
func (r *Relay) Dispatch(path string, payload map[string]string) []error {
	r.mu.RLock()
	names := make([]string, 0, len(r.handlers))
	copy := make(map[string]Handler, len(r.handlers))
	for n, h := range r.handlers {
		names = append(names, n)
		copy[n] = h
	}
	r.mu.RUnlock()

	var errs []error
	for _, n := range names {
		if err := copy[n](path, payload); err != nil {
			errs = append(errs, fmt.Errorf("secretrelay: handler %q: %w", n, err))
		}
	}
	return errs
}

// Len returns the number of registered handlers.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.handlers)
}
