// Package renew provides automatic lease and token renewal tracking.
package renew

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrAlreadyTracked is returned when a path is already being renewed.
var ErrAlreadyTracked = errors.New("renew: path already tracked")

// Entry holds renewal metadata for a single secret path.
type Entry struct {
	Path      string
	LeaseTTL  time.Duration
	RenewAt   time.Time
	RenewedAt time.Time
	RenewCount int
}

// RenewFunc is called when a secret is due for renewal.
type RenewFunc func(ctx context.Context, path string) (time.Duration, error)

// Renewer tracks secrets and triggers renewal before expiry.
type Renewer struct {
	mu      sync.Mutex
	entries map[string]*Entry
	thresh  float64 // fraction of TTL at which to renew (e.g. 0.75)
	now     func() time.Time
}

// New creates a Renewer that renews secrets when thresh fraction of TTL has elapsed.
func New(renewThreshold float64) (*Renewer, error) {
	if renewThreshold <= 0 || renewThreshold >= 1 {
		return nil, fmt.Errorf("renew: threshold must be between 0 and 1, got %v", renewThreshold)
	}
	return &Renewer{
		entries: make(map[string]*Entry),
		thresh:  renewThreshold,
		now:     time.Now,
	}, nil
}

// Register adds a path to renewal tracking.
func (r *Renewer) Register(path string, leaseTTL time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[path]; ok {
		return ErrAlreadyTracked
	}
	r.entries[path] = &Entry{
		Path:     path,
		LeaseTTL: leaseTTL,
		RenewAt:  r.now().Add(time.Duration(float64(leaseTTL) * r.thresh)),
	}
	return nil
}

// Due returns all entries whose renewal time has passed.
func (r *Renewer) Due() []*Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	var out []*Entry
	for _, e := range r.entries {
		if !now.Before(e.RenewAt) {
			out = append(out, e)
		}
	}
	return out
}

// RecordRenewal updates an entry after a successful renewal.
func (r *Renewer) RecordRenewal(path string, newTTL time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[path]
	if !ok {
		return fmt.Errorf("renew: path not tracked: %s", path)
	}
	now := r.now()
	e.LeaseTTL = newTTL
	e.RenewedAt = now
	e.RenewCount++
	e.RenewAt = now.Add(time.Duration(float64(newTTL) * r.thresh))
	return nil
}

// Remove stops tracking a path.
func (r *Renewer) Remove(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, path)
}
