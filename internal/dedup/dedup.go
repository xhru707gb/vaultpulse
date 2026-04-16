// Package dedup provides alert deduplication to suppress repeated notifications
// for the same secret path and event type within a configurable window.
package dedup

import (
	"sync"
	"time"
)

// Key uniquely identifies an alert event.
type Key struct {
	Path  string
	Event string
}

// Deduplicator suppresses duplicate alerts within a time window.
type Deduplicator struct {
	mu     sync.Mutex
	seen   map[Key]time.Time
	window time.Duration
	now    func() time.Time
}

// New creates a Deduplicator with the given suppression window.
func New(window time.Duration) (*Deduplicator, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Deduplicator{
		seen:   make(map[Key]time.Time),
		window: window,
		now:    time.Now,
	}, nil
}

// IsDuplicate returns true if the key was seen within the suppression window.
// If not a duplicate, it records the key and returns false.
func (d *Deduplicator) IsDuplicate(path, event string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	k := Key{Path: path, Event: event}
	if t, ok := d.seen[k]; ok && d.now().Sub(t) < d.window {
		return true
	}
	d.seen[k] = d.now()
	return false
}

// Flush removes all entries older than the suppression window.
func (d *Deduplicator) Flush() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	cutoff := d.now().Add(-d.window)
	removed := 0
	for k, t := range d.seen {
		if t.Before(cutoff) {
			delete(d.seen, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of tracked keys.
func (d *Deduplicator) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.seen)
}
