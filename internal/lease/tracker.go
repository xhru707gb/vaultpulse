// Package lease tracks Vault lease IDs and their renewal deadlines.
package lease

import (
	"errors"
	"sync"
	"time"
)

// ErrNotFound is returned when a lease ID is not registered.
var ErrNotFound = errors.New("lease: not found")

// Entry holds metadata for a single Vault lease.
type Entry struct {
	LeaseID   string
	Path      string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Renewable bool
}

// TTL returns the remaining duration until expiry.
func (e Entry) TTL(now time.Time) time.Duration {
	return e.ExpiresAt.Sub(now)
}

// Tracker stores and queries active lease entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns a new Tracker.
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{entries: make(map[string]Entry), now: now}
}

// Register adds or replaces a lease entry.
func (t *Tracker) Register(e Entry) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[e.LeaseID] = e
}

// Remove deletes a lease by ID.
func (t *Tracker) Remove(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[id]; !ok {
		return ErrNotFound
	}
	delete(t.entries, id)
	return nil
}

// Get returns the entry for a lease ID.
func (t *Tracker) Get(id string) (Entry, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[id]
	if !ok {
		return Entry{}, ErrNotFound
	}
	return e, nil
}

// Expiring returns all leases expiring within the given threshold.
func (t *Tracker) Expiring(threshold time.Duration) []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	now := t.now()
	var out []Entry
	for _, e := range t.entries {
		if e.ExpiresAt.Sub(now) <= threshold {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the number of tracked leases.
func (t *Tracker) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.entries)
}
