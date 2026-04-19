// Package semaphore provides a counting semaphore for bounding concurrent
// Vault API calls within vaultpulse.
package semaphore

import (
	"context"
	"errors"
	"sync"
)

// ErrInvalidSize is returned when a non-positive size is supplied.
var ErrInvalidSize = errors.New("semaphore: size must be greater than zero")

// Semaphore is a counting semaphore backed by a buffered channel.
type Semaphore struct {
	mu      sync.Mutex
	slots   chan struct{}
	size    int
	acquired int
}

// New creates a Semaphore with the given capacity.
func New(size int) (*Semaphore, error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}
	return &Semaphore{
		slots: make(chan struct{}, size),
		size:  size,
	}, nil
}

// Acquire blocks until a slot is available or ctx is cancelled.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		s.mu.Lock()
		s.acquired++
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release frees a previously acquired slot.
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
		s.mu.Lock()
		s.acquired--
		s.mu.Unlock()
	default:
	}
}

// Acquired returns the number of currently held slots.
func (s *Semaphore) Acquired() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.acquired
}

// Size returns the total capacity of the semaphore.
func (s *Semaphore) Size() int {
	return s.size
}
