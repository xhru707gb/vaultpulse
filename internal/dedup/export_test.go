package dedup

// SetNow overrides the clock for testing.
func (d *Deduplicator) SetNow(fn func() time.Time) {
	d.now = fn
}

import "time"
