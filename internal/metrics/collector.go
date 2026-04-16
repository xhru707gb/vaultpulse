// Package metrics records per-path check snapshots and exposes
// lightweight aggregation for CLI reporting.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds the result of a single check run for one secret path.
type Snapshot struct {
	Path      string
	CheckedAt time.Time
	Duration  time.Duration
	Status    string // "ok", "warning", "expired"
	Error     string // empty when no error
}

// Collector stores the most-recent snapshot per path.
type Collector struct {
	mu   sync.RWMutex
	data map[string]Snapshot
	now  func() time.Time
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{
		data: make(map[string]Snapshot),
		now:  time.Now,
	}
}

// Record stores (or overwrites) the snapshot for path.
func (c *Collector) Record(path, status string, d time.Duration, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	c.data[path] = Snapshot{
		Path:      path,
		CheckedAt: c.now(),
		Duration:  d,
		Status:    status,
		Error:     errStr,
	}
}

// Get returns the snapshot for path and whether it exists.
func (c *Collector) Get(path string) (Snapshot, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.data[path]
	return s, ok
}

// All returns a copy of every stored snapshot.
func (c *Collector) All() []Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]Snapshot, 0, len(c.data))
	for _, s := range c.data {
		out = append(out, s)
	}
	return out
}

// Reset removes all stored snapshots.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]Snapshot)
}
