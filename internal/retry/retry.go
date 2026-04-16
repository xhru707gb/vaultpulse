// Package retry provides a simple exponential backoff retry mechanism
// for transient errors encountered when communicating with Vault.
package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Config holds retry behaviour parameters.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// BaseDelay is the initial wait duration before the first retry.
	BaseDelay time.Duration
	// MaxDelay caps the computed exponential delay.
	MaxDelay time.Duration
}

// DefaultConfig returns sensible defaults for production use.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 4,
		BaseDelay:   250 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// ErrMaxAttemptsReached is returned when all attempts are exhausted.
var ErrMaxAttemptsReached = errors.New("retry: max attempts reached")

// Do executes fn up to cfg.MaxAttempts times, backing off exponentially
// between attempts. It respects ctx cancellation. The last non-nil error
// returned by fn is wrapped and returned alongside ErrMaxAttemptsReached.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := fn(); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt == cfg.MaxAttempts-1 {
			break
		}

		delay := delay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}

	return errors.Join(ErrMaxAttemptsReached, lastErr)
}

// delay computes the exponential backoff for the given attempt index.
func delay(base, max time.Duration, attempt int) time.Duration {
	mult := math.Pow(2, float64(attempt))
	d := time.Duration(float64(base) * mult)
	if d > max {
		return max
	}
	return d
}
