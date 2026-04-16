// Package throttle provides a simple concurrent request throttler
// that limits the number of in-flight Vault API calls at any time.
package throttle

import (
	"context"
	"errors"
)

// ErrThrottled is returned when the throttle is at capacity and the context expires.
var ErrThrottled = errors.New("throttle: request rejected, at capacity")

// Throttle limits concurrent operations using a semaphore channel.
type Throttle struct {
	sem chan struct{}
}

// New creates a Throttle allowing at most concurrency simultaneous operations.
// It returns an error if concurrency is less than 1.
func New(concurrency int) (*Throttle, error) {
	if concurrency < 1 {
		return nil, errors.New("throttle: concurrency must be at least 1")
	}
	return &Throttle{sem: make(chan struct{}, concurrency)}, nil
}

// Acquire blocks until a slot is available or ctx is done.
// Returns ErrThrottled if the context expires before a slot is acquired.
func (t *Throttle) Acquire(ctx context.Context) error {
	select {
	case t.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ErrThrottled
	}
}

// Release frees a previously acquired slot.
func (t *Throttle) Release() {
	<-t.sem
}

// Do acquires a slot, runs fn, then releases the slot.
// Returns ErrThrottled if ctx expires before acquisition.
func (t *Throttle) Do(ctx context.Context, fn func() error) error {
	if err := t.Acquire(ctx); err != nil {
		return err
	}
	defer t.Release()
	return fn()
}

// Capacity returns the maximum number of concurrent operations.
func (t *Throttle) Capacity() int {
	return cap(t.sem)
}

// InFlight returns the current number of in-flight operations.
func (t *Throttle) InFlight() int {
	return len(t.sem)
}
