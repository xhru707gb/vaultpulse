// Package secretreview provides periodic review scheduling and tracking
// for secrets that require manual or automated sign-off.
package secretreview

import (
	"errors"
	"sync"
	"time"
)

// Status represents the review state of a secret.
type Status int

const (
	StatusPending Status = iota
	StatusApproved
	StatusOverdue
)

// Entry holds review metadata for a single secret path.
type Entry struct {
	Path        string
	Interval    time.Duration
	LastReview  time.Time
	NextReview  time.Time
	Reviewer    string
	Status      Status
}

// Reviewer tracks secret review schedules.
type Reviewer struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	now     func() time.Time
}

// New creates a Reviewer with an optional clock override.
func New(now func() time.Time) (*Reviewer, error) {
	if now == nil {
		now = time.Now
	}
	return &Reviewer{
		entries: make(map[string]*Entry),
		now:     now,
	}, nil
}

// Register adds or updates a secret review entry.
func (r *Reviewer) Register(path, reviewer string, interval time.Duration, lastReview time.Time) error {
	if path == "" {
		return errors.New("path must not be empty")
	}
	if reviewer == "" {
		return errors.New("reviewer must not be empty")
	}
	if interval <= 0 {
		return errors.New("interval must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[path] = &Entry{
		Path:       path,
		Interval:   interval,
		LastReview: lastReview,
		NextReview: lastReview.Add(interval),
		Reviewer:   reviewer,
	}
	return nil
}

// Evaluate returns all entries with computed statuses.
func (r *Reviewer) Evaluate() []*Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()
	now := r.now()
	out := make([]*Entry, 0, len(r.entries))
	for _, e := range r.entries {
		copy := *e
		if now.Before(copy.NextReview) {
			copy.Status = StatusApproved
		} else {
			copy.Status = StatusOverdue
		}
		out = append(out, &copy)
	}
	return out
}

// Approve marks a secret as reviewed now, resetting the next review time.
func (r *Reviewer) Approve(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.entries[path]
	if !ok {
		return errors.New("unknown path: " + path)
	}
	now := r.now()
	e.LastReview = now
	e.NextReview = now.Add(e.Interval)
	e.Status = StatusApproved
	return nil
}
