package circuit

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
