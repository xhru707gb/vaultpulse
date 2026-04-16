// Package backoff provides configurable exponential back-off with jitter.
package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

// Config holds back-off parameters.
type Config struct {
	InitialInterval time.Duration
	Multiplier      float64
	MaxInterval     time.Duration
	Jitter          bool
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		InitialInterval: 200 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     30 * time.Second,
		Jitter:          true,
	}
}

// Backoff computes successive wait durations.
type Backoff struct {
	cfg     Config
	attempt int
}

// New validates cfg and returns a Backoff.
func New(cfg Config) (*Backoff, error) {
	if cfg.InitialInterval <= 0 {
		return nil, errors.New("backoff: InitialInterval must be positive")
	}
	if cfg.Multiplier < 1 {
		return nil, errors.New("backoff: Multiplier must be >= 1")
	}
	if cfg.MaxInterval < cfg.InitialInterval {
		return nil, errors.New("backoff: MaxInterval must be >= InitialInterval")
	}
	return &Backoff{cfg: cfg}, nil
}

// Next returns the duration to wait before the next attempt and increments
// the internal attempt counter.
func (b *Backoff) Next() time.Duration {
	raw := float64(b.cfg.InitialInterval) * math.Pow(b.cfg.Multiplier, float64(b.attempt))
	if raw > float64(b.cfg.MaxInterval) {
		raw = float64(b.cfg.MaxInterval)
	}
	b.attempt++
	if b.cfg.Jitter {
		raw = raw/2 + rand.Float64()*raw/2
	}
	return time.Duration(raw)
}

// Reset sets the attempt counter back to zero.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt index (zero-based).
func (b *Backoff) Attempt() int { return b.attempt }
