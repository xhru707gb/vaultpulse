// Package presign provides URL pre-signing tracking for Vault secrets,
// recording when a signed URL was issued and when it expires.
package presign

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a pre-signed URL record.
type Entry struct {
	Path      string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Token     string
}

// IsExpired reports whether the entry has passed its expiry time.
func (e Entry) IsExpired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

// TTL returns the remaining duration until expiry.
func (e Entry) TTL(now time.Time) time.Duration {
	if e.IsExpired(now) {
		return 0
	}
	return e.ExpiresAt.Sub(now)
}

// Tracker stores and manages pre-signed URL entries.
type Tracker struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Register records a new pre-signed entry for path.
func (t *Tracker) Register(path, token string, ttl time.Duration) error {
	if path == "" {
		return errors.New("presign: path must not be empty")
	}
	if token == "" {
		return errors.New("presign: token must not be empty")
	}
	if ttl <= 0 {
		return errors.New("presign: ttl must be positive")
	}
	now := t.now()
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[path] = Entry{
		Path:      path,
		IssuedAt:  now,
		ExpiresAt: now.Add(ttl),
		Token:     token,
	}
	return nil
}

// Get returns the entry for path, if present.
func (t *Tracker) Get(path string) (Entry, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[path]
	return e, ok
}

// All returns all tracked entries.
func (t *Tracker) All() []Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// Revoke removes the entry for path.
func (t *Tracker) Revoke(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[path]; !ok {
		return errors.New("presign: path not found")
	}
	delete(t.entries, path)
	return nil
}
