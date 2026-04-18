// Package jitter provides utilities for adding randomised jitter to
// durations, useful for spreading out retry or polling intervals.
package jitter

import (
	"errors"
	"math/rand"
	"time"
)

// ErrInvalidFactor is returned when the jitter factor is out of range.
var ErrInvalidFactor = errors.New("jitter: factor must be between 0 and 1")

// Config holds jitter configuration.
type Config struct {
	// Factor controls the maximum fraction of base that may be added.
	// Must be in the range [0, 1].
	Factor float64
	// Source is the random source to use. If nil, a default source is used.
	Source rand.Source
}

// Jitter applies randomised jitter to durations.
type Jitter struct {
	cfg Config
	rng *rand.Rand
}

// New creates a new Jitter. Returns ErrInvalidFactor if cfg.Factor is
// outside [0, 1].
func New(cfg Config) (*Jitter, error) {
	if cfg.Factor < 0 || cfg.Factor > 1 {
		return nil, ErrInvalidFactor
	}
	src := cfg.Source
	if src == nil {
		src = rand.NewSource(time.Now().UnixNano())
	}
	return &Jitter{cfg: cfg, rng: rand.New(src)}, nil
}

// Apply returns base plus a random duration in [0, base*Factor].
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if j.cfg.Factor == 0 || base <= 0 {
		return base
	}
	max := float64(base) * j.cfg.Factor
	delta := time.Duration(j.rng.Float64() * max)
	return base + delta
}

// ApplyRange returns a duration uniformly distributed in
// [base*(1-Factor), base*(1+Factor)].
func (j *Jitter) ApplyRange(base time.Duration) time.Duration {
	if j.cfg.Factor == 0 || base <= 0 {
		return base
	}
	spread := float64(base) * j.cfg.Factor
	delta := time.Duration((j.rng.Float64()*2-1)*spread)
	return base + delta
}
