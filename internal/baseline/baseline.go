// Package baseline records and compares secret metadata snapshots to detect drift.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry holds a point-in-time record of a secret's metadata.
type Entry struct {
	Path        string    `json:"path"`
	Version     int       `json:"version"`
	TTL         int64     `json:"ttl_seconds"`
	LastRotated time.Time `json:"last_rotated"`
	CapturedAt  time.Time `json:"captured_at"`
}

// Drift describes a deviation from the recorded baseline.
type Drift struct {
	Path  string
	Field string
	Was   string
	Now   string
}

// Store holds baseline entries keyed by secret path.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
	now     func() time.Time
}

// New returns an empty baseline Store.
func New() *Store {
	return &Store{
		entries: make(map[string]Entry),
		now:     time.Now,
	}
}

// Record saves or overwrites the baseline entry for a path.
func (s *Store) Record(e Entry) {
	e.CapturedAt = s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[e.Path] = e
}

// Compare returns drifts between the stored baseline and the supplied entry.
// Returns an error if no baseline exists for the path.
func (s *Store) Compare(current Entry) ([]Drift, error) {
	s.mu.RLock()
	base, ok := s.entries[current.Path]
	s.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("baseline: no entry for path %q", current.Path)
	}
	var drifts []Drift
	if base.Version != current.Version {
		drifts = append(drifts, Drift{Path: current.Path, Field: "version",
			Was: fmt.Sprintf("%d", base.Version), Now: fmt.Sprintf("%d", current.Version)})
	}
	if base.TTL != current.TTL {
		drifts = append(drifts, Drift{Path: current.Path, Field: "ttl_seconds",
			Was: fmt.Sprintf("%d", base.TTL), Now: fmt.Sprintf("%d", current.TTL)})
	}
	if !base.LastRotated.Equal(current.LastRotated) {
		drifts = append(drifts, Drift{Path: current.Path, Field: "last_rotated",
			Was: base.LastRotated.UTC().Format(time.RFC3339),
			Now: current.LastRotated.UTC().Format(time.RFC3339)})
	}
	return drifts, nil
}

// SaveJSON persists all entries to a JSON file.
func (s *Store) SaveJSON(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s.entries)
}

// LoadJSON restores entries from a JSON file.
func (s *Store) LoadJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.NewDecoder(f).Decode(&s.entries)
}
