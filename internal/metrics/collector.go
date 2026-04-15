// Package metrics provides aggregated runtime metrics collection
// across expiry, rotation, health, and policy check results.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	CollectedAt     time.Time
	TotalSecrets    int
	Expired         int
	Warning         int
	Healthy         int
	OverdueRotation int
	PolicyViolation int
	LastCheckDur    time.Duration
}

// Collector accumulates metrics from check runs.
type Collector struct {
	mu       sync.Mutex
	current  Snapshot
	nowFn    func() time.Time
}

// NewCollector returns a new Collector using real wall-clock time.
func NewCollector() *Collector {
	return &Collector{nowFn: time.Now}
}

// Record replaces the current snapshot with the provided values and
// stamps the collection time.
func (c *Collector) Record(s Snapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s.CollectedAt = c.nowFn()
	c.current = s
}

// Get returns a copy of the most recently recorded snapshot.
func (c *Collector) Get() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.current
}

// Reset zeroes the current snapshot.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = Snapshot{}
}
