// Package circuit implements a circuit breaker to protect against
// repeated failures when communicating with external services such as Vault.
package circuit

import (
	"errors"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// ErrOpen is returned when the circuit is open.
var ErrOpen = errors.New("circuit breaker is open")

// Config holds circuit breaker configuration.
type Config struct {
	MaxFailures int
	OpenTimeout time.Duration
}

// Breaker is a simple circuit breaker.
type Breaker struct {
	mu          sync.Mutex
	cfg         Config
	failures    int
	state       State
	openedAt    time.Time
	now         func() time.Time
}

// New creates a new Breaker with the given config.
func New(cfg Config) (*Breaker, error) {
	if cfg.MaxFailures <= 0 {
		return nil, errors.New("MaxFailures must be > 0")
	}
	if cfg.OpenTimeout <= 0 {
		return nil, errors.New("OpenTimeout must be > 0")
	}
	return &Breaker{cfg: cfg, now: time.Now}, nil
}

// Allow returns nil if the call is permitted, or ErrOpen if the circuit is open.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.cfg.OpenTimeout {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	}
	return nil
}

// RecordSuccess resets the breaker on success.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.cfg.MaxFailures {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current state of the breaker.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
