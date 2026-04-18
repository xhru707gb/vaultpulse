package window

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a timestamped value recorded in the window.
type Entry[T any] struct {
	Value     T
	RecordedAt time.Time
}

// Window is a sliding time window that retains entries within a duration.
type Window[T any] struct {
	mu       sync.Mutex
	duration time.Duration
	entries  []Entry[T]
	now      func() time.Time
}

// New creates a Window with the given retention duration.
func New[T any](duration time.Duration) (*Window[T], error) {
	if duration <= 0 {
		return nil, errors.New("window: duration must be positive")
	}
	return &Window[T]{duration: duration, now: time.Now}, nil
}

// Add inserts a value into the window, pruning stale entries first.
func (w *Window[T]) Add(v T) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	w.entries = append(w.entries, Entry[T]{Value: v, RecordedAt: w.now()})
}

// Entries returns a copy of all current (non-stale) entries.
func (w *Window[T]) Entries() []Entry[T] {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	out := make([]Entry[T], len(w.entries))
	copy(out, w.entries)
	return out
}

// Len returns the number of active entries.
func (w *Window[T]) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.prune()
	return len(w.entries)
}

// Reset clears all entries.
func (w *Window[T]) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = nil
}

func (w *Window[T]) prune() {
	cutoff := w.now().Add(-w.duration)
	i := 0
	for i < len(w.entries) && w.entries[i].RecordedAt.Before(cutoff) {
		i++
	}
	w.entries = w.entries[i:]
}
