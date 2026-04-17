package ratelimit

import (
	"context"
	"errors"
	"fmt"
)

// ErrRateLimited is returned when a request is denied by the limiter.
var ErrRateLimited = errors.New("rate limit exceeded")

// Do executes fn only if the limiter allows the request.
// If the token bucket is exhausted, ErrRateLimited is returned immediately.
func Do(ctx context.Context, l *Limiter, fn func(ctx context.Context) error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !l.Allow() {
		return fmt.Errorf("%w: retry after %.1fs",
			ErrRateLimited, l.RetryAfter().Seconds())
	}

	return fn(ctx)
}

// DoWithKey executes fn only if the per-key limiter allows the request.
// A separate Limiter is maintained per key using the provided map.
func DoWithKey(
	ctx context.Context,
	key string,
	limiters map[string]*Limiter,
	cfg Config,
	fn func(ctx context.Context) error,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	l, ok := limiters[key]
	if !ok {
		var err error
		l, err = New(cfg)
		if err != nil {
			return fmt.Errorf("ratelimit: create limiter for key %q: %w", key, err)
		}
		limiters[key] = l
	}

	if !l.Allow() {
		return fmt.Errorf("%w: key=%s retry after %.1fs",
			ErrRateLimited, key, l.RetryAfter().Seconds())
	}

	return fn(ctx)
}
