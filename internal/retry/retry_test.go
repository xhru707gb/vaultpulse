package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/retry"
)

var errTransient = errors.New("transient error")

func fastConfig(maxAttempts int) retry.Config {
	return retry.Config{
		MaxAttempts: maxAttempts,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    5 * time.Millisecond,
	}
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(3), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUntilSuccess(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(4), func() error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttemptsReturnsError(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), fastConfig(3), func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected wrapped errTransient, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ContextCancelledStopsRetries(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := retry.Do(ctx, fastConfig(10), func() error {
		calls++
		cancel()
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call before cancel, got %d", calls)
	}
}

func TestDo_ZeroMaxAttemptsRunsOnce(t *testing.T) {
	calls := 0
	err := retry.Do(context.Background(), retry.Config{MaxAttempts: 0, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}, func() error {
		calls++
		return errTransient
	})
	if !errors.Is(err, retry.ErrMaxAttemptsReached) {
		t.Fatalf("expected ErrMaxAttemptsReached, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls)
	}
}
