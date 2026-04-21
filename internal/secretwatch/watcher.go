// Package secretwatch provides a periodic watcher that monitors secret
// metadata changes and emits events when secrets are added, removed, or
// modified between polling cycles.
package secretwatch

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Event describes a change detected during a watch cycle.
type Event struct {
	Path   string
	Kind   string // "added", "removed", "modified"
	Detail string
}

// Handler is called with the list of events detected each cycle.
type Handler func(events []Event)

// SecretLister returns the current set of secret paths and their versions.
type SecretLister func(ctx context.Context) (map[string]int, error)

// Watcher polls a SecretLister on a fixed interval and invokes a Handler
// whenever the observed state differs from the previous cycle.
type Watcher struct {
	lister   SecretLister
	handler  Handler
	interval time.Duration
	mu       sync.Mutex
	prev     map[string]int
}

// New creates a Watcher. interval must be positive; lister and handler must
// be non-nil.
func New(lister SecretLister, handler Handler, interval time.Duration) (*Watcher, error) {
	if lister == nil {
		return nil, errors.New("secretwatch: lister must not be nil")
	}
	if handler == nil {
		return nil, errors.New("secretwatch: handler must not be nil")
	}
	if interval <= 0 {
		return nil, errors.New("secretwatch: interval must be positive")
	}
	return &Watcher{
		lister:   lister,
		handler:  handler,
		interval: interval,
		prev:     make(map[string]int),
	}, nil
}

// Run starts the watch loop. It blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *Watcher) poll(ctx context.Context) {
	current, err := w.lister(ctx)
	if err != nil {
		return
	}
	w.mu.Lock()
	events := diff(w.prev, current)
	w.prev = current
	w.mu.Unlock()
	if len(events) > 0 {
		w.handler(events)
	}
}

func diff(prev, current map[string]int) []Event {
	var events []Event
	for path, ver := range current {
		if oldVer, ok := prev[path]; !ok {
			events = append(events, Event{Path: path, Kind: "added"})
		} else if ver != oldVer {
			events = append(events, Event{Path: path, Kind: "modified"})
		}
	}
	for path := range prev {
		if _, ok := current[path]; !ok {
			events = append(events, Event{Path: path, Kind: "removed"})
		}
	}
	return events
}
