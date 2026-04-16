// Package watch provides periodic polling of Vault secrets,
// triggering callbacks when expiry or rotation state changes.
package watch

import (
	"context"
	"time"

	"github.com/your-org/vaultpulse/internal/expiry"
)

// Event holds the result of a single poll cycle.
type Event struct {
	Statuses []expiry.Status
	Err      error
	At       time.Time
}

// Handler is called after every poll cycle.
type Handler func(Event)

// Watcher polls the expiry checker at a fixed interval.
type Watcher struct {
	checker  *expiry.Checker
	interval time.Duration
	handler  Handler
}

// New creates a Watcher that polls checker every interval and calls handler.
func New(checker *expiry.Checker, interval time.Duration, handler Handler) (*Watcher, error) {
	if interval <= 0 {
		return nil, ErrInvalidInterval
	}
	if handler == nil {
		return nil, ErrNilHandler
	}
	return &Watcher{checker: checker, interval: interval, handler: handler}, nil
}

// Run starts the polling loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context, paths []string) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			statuses, err := w.checker.CheckAll(ctx, paths)
			w.handler(Event{Statuses: statuses, Err: err, At: t})
		}
	}
}
