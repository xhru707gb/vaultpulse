package cooldown

// SetNow overrides the time source for testing.
func (t *Tracker) SetNow(fn func() time.Time) {
	t.now = fn
}
