package circuit

import "time"

// SetNow overrides the time source for testing.
func (b *Breaker) SetNow(fn func() time.Time) {
	b.now = fn
}

// Failures exposes the internal failure counter for testing.
func (b *Breaker) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}

// State exposes the internal circuit state for testing.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
