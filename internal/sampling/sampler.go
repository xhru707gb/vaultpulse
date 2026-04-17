// Package sampling provides probabilistic and rate-based sampling for audit events.
package sampling

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// ErrInvalidRate is returned when the sample rate is out of range.
var ErrInvalidRate = errors.New("sampling: rate must be between 0.0 and 1.0")

// Config holds sampler configuration.
type Config struct {
	// Rate is the probability [0.0, 1.0] that a given event is sampled.
	Rate float64
	// Seed is used to initialise the RNG; 0 uses a time-based seed.
	Seed int64
}

// Sampler decides whether an event should be forwarded.
type Sampler struct {
	mu   sync.Mutex
	rng  *rand.Rand
	rate float64
}

// New creates a Sampler from cfg.
func New(cfg Config) (*Sampler, error) {
	if cfg.Rate < 0.0 || cfg.Rate > 1.0 {
		return nil, ErrInvalidRate
	}
	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &Sampler{
		rng:  rand.New(rand.NewSource(seed)), //nolint:gosec
		rate: cfg.Rate,
	}, nil
}

// Sample returns true if the event should be kept.
func (s *Sampler) Sample() bool {
	if s.rate >= 1.0 {
		return true
	}
	if s.rate <= 0.0 {
		return false
	}
	s.mu.Lock()
	v := s.rng.Float64()
	s.mu.Unlock()
	return v < s.rate
}

// Rate returns the configured sample rate.
func (s *Sampler) Rate() float64 { return s.rate }
