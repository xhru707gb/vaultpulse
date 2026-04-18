package window

// SetNow overrides the clock used by the window for testing.
func (w *Window[T]) SetNow(fn func() time.Time) {
	w.now = fn
}
