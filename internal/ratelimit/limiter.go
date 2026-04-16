// Package ratelimit provides a simple token-bucket rate limiter for
// controlling the frequency of Vault API requests made by vaultpulse.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limiter enforces a maximum number of requests per second.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	rate     float64 // tokens per nanosecond
	lastTick time.Time
	now      func() time.Time
}

// Config holds rate limiter configuration.
type Config struct {
	// RequestsPerSecond is the sustained request rate allowed.
	RequestsPerSecond float64
	// Burst is the maximum number of requests allowed in a single instant.
	Burst float64
}

// New creates a Limiter from the given Config.
// Returns an error if RequestsPerSecond or Burst are non-positive.
func New(cfg Config) (*Limiter, error) {
	if cfg.RequestsPerSecond <= 0 {
		return nil, fmt.Errorf("ratelimit: RequestsPerSecond must be positive, got %v", cfg.RequestsPerSecond)
	}
	if cfg.Burst <= 0 {
		return nil, fmt.Errorf("ratelimit: Burst must be positive, got %v", cfg.Burst)
	}
	return &Limiter{
		tokens:   cfg.Burst,
		max:      cfg.Burst,
		rate:     cfg.RequestsPerSecond / float64(time.Second),
		lastTick: time.Now(),
		now:      time.Now,
	}, nil
}

// Allow reports whether a single request may proceed.
// It refills tokens based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.lastTick)
	l.lastTick = now

	l.tokens += float64(elapsed) * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// Wait blocks until a token is available or the context deadline is exceeded.
func (l *Limiter) Wait() {
	for {
		if l.Allow() {
			return
		}
		time.Sleep(time.Millisecond * 5)
	}
}

// Tokens returns the current number of available tokens without consuming any.
func (l *Limiter) Tokens() float64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	elapsed := now.Sub(l.lastTick)

	tokens := l.tokens + float64(elapsed)*l.rate
	if tokens > l.max {
		tokens = l.max
	}
	return tokens
}
