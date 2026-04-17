package schedule_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vaultpulse/internal/schedule"
)

func TestNew_InvalidInterval(t *testing.T) {
	_, err := schedule.New(0, func(ctx context.Context) error { return nil })
	if !errors.Is(err, schedule.ErrInvalidInterval) {
		t.Fatalf("expected ErrInvalidInterval, got %v", err)
	}
}

func TestNew_NilJob(t *testing.T) {
	_, err := schedule.New(time.Second, nil)
	if err == nil {
		t.Fatal("expected error for nil job")
	}
}

func TestNew_Valid(t *testing.T) {
	s, err := schedule.New(time.Second, func(ctx context.Context) error { return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestRun_ExecutesJobOnTick(t *testing.T) {
	var count atomic.Int32
	s, _ := schedule.New(20*time.Millisecond, func(ctx context.Context) error {
		count.Add(1)
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	_ = s.Run(ctx)
	if count.Load() < 2 {
		t.Fatalf("expected at least 2 executions, got %d", count.Load())
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	s, _ := schedule.New(10*time.Millisecond, func(ctx context.Context) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := s.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestRun_StopsOnJobError(t *testing.T) {
	sentinel := errors.New("job failed")
	s, _ := schedule.New(10*time.Millisecond, func(ctx context.Context) error {
		return sentinel
	})
	ctx := context.Background()
	err := s.Run(ctx)
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}
