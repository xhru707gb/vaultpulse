// Package rollup aggregates multiple secret status events into a single
// summarised payload, reducing noise in downstream alert channels.
package rollup

import (
	"sync"
	"time"
)

// Event represents a single status event to be rolled up.
type Event struct {
	Path      string
	Level     string // "ok", "warning", "expired"
	Message   string
	OccuredAt time.Time
}

// Summary is the aggregated output produced by Flush.
type Summary struct {
	Window    time.Duration
	Total     int
	ByLevel   map[string]int
	Events    []Event
	FlushedAt time.Time
}

// Aggregator buffers events within a time window and flushes them as a Summary.
type Aggregator struct {
	mu     sync.Mutex
	window time.Duration
	events []Event
	now    func() time.Time
}

// New creates an Aggregator with the given window duration.
// Returns an error if window is non-positive.
func New(window time.Duration) (*Aggregator, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Aggregator{window: window, now: time.Now}, nil
}

// Add appends an event to the internal buffer.
func (a *Aggregator) Add(e Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if e.OccuredAt.IsZero() {
		e.OccuredAt = a.now()
	}
	a.events = append(a.events, e)
}

// Flush drains the buffer and returns a Summary. The buffer is cleared.
func (a *Aggregator) Flush() Summary {
	a.mu.Lock()
	defer a.mu.Unlock()

	byLevel := make(map[string]int)
	for _, e := range a.events {
		byLevel[e.Level]++
	}

	s := Summary{
		Window:    a.window,
		Total:     len(a.events),
		ByLevel:   byLevel,
		Events:    append([]Event(nil), a.events...),
		FlushedAt: a.now(),
	}
	a.events = a.events[:0]
	return s
}

// Len returns the number of buffered events.
func (a *Aggregator) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.events)
}
