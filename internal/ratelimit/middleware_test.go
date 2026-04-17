package ratelimit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/ratelimit"
)

func testCfg(rate float64, burst int) ratelimit.Config {
	return ratelimit.Config{
		Rate:     rate,
		Burst:    burst,
		Interval: time.Second,
	}
}

func TestDo_AllowsWhenTokensAvailable(t *testing.T) {
	l, _ := ratelimit.New(testCfg(10, 5))
	called := false
	err := ratelimit.Do(context.Background(), l, func(_ context.Context) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected fn to be called")
	}
}

func TestDo_RateLimitedWhenExhausted(t *testing.T) {
	l, _ := ratelimit.New(testCfg(0.001, 1))
	// consume the single burst token
	ratelimit.Do(context.Background(), l, func(_ context.Context) error { return nil }) //nolint

	err := ratelimit.Do(context.Background(), l, func(_ context.Context) error {
		return nil
	})
	if !errors.Is(err, ratelimit.ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
}

func TestDo_CancelledContext(t *testing.T) {
	l, _ := ratelimit.New(testCfg(10, 5))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ratelimit.Do(ctx, l, func(_ context.Context) error { return nil })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestDoWithKey_CreatesLimiterPerKey(t *testing.T) {
	limiters := make(map[string]*ratelimit.Limiter)
	cfg := testCfg(10, 5)
	called := 0
	for i := 0; i < 3; i++ {
		err := ratelimit.DoWithKey(context.Background(), "secret/a", limiters, cfg,
			func(_ context.Context) error { called++; return nil })
		if err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}
	if called != 3 {
		t.Fatalf("expected 3 calls, got %d", called)
	}
	if len(limiters) != 1 {
		t.Fatalf("expected 1 limiter entry, got %d", len(limiters))
	}
}

func TestDoWithKey_IndependentLimitsPerKey(t *testing.T) {
	limiters := make(map[string]*ratelimit.Limiter)
	cfg := testCfg(0.001, 1)

	// exhaust key a
	ratelimit.DoWithKey(context.Background(), "a", limiters, cfg, func(_ context.Context) error { return nil }) //nolint

	// key b should still be allowed
	err := ratelimit.DoWithKey(context.Background(), "b", limiters, cfg,
		func(_ context.Context) error { return nil })
	if err != nil {
		t.Fatalf("key b should not be rate limited, got %v", err)
	}
}
