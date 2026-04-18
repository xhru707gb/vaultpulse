package tokenwatch

import "time"

// NewWithClock creates a Watcher with an injectable clock for testing.
func NewWithClock(f Fetcher, warnThreshold time.Duration, now func() time.Time) (*Watcher, error) {
	w, err := New(f, warnThreshold)
	if err != nil {
		return nil, err
	}
	w.now = now
	return w, nil
}
