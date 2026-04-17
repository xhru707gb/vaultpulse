package schedule_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vaultpulse/internal/schedule"
)

func TestRegister_Valid(t *testing.T) {
	r := schedule.NewRunner()
	s, _ := schedule.New(time.Second, func(ctx context.Context) error { return nil })
	if err := r.Register("check", s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_EmptyName(t *testing.T) {
	r := schedule.NewRunner()
	s, _ := schedule.New(time.Second, func(ctx context.Context) error { return nil })
	if err := r.Register("", s); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilScheduler(t *testing.T) {
	r := schedule.NewRunner()
	if err := r.Register("check", nil); err == nil {
		t.Fatal("expected error for nil scheduler")
	}
}

func TestRegister_Duplicate(t *testing.T) {
	r := schedule.NewRunner()
	s, _ := schedule.New(time.Second, func(ctx context.Context) error { return nil })
	_ = r.Register("check", s)
	if err := r.Register("check", s); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestRunAll_StopsOnContextCancel(t *testing.T) {
	r := schedule.NewRunner()
	s, _ := schedule.New(10*time.Millisecond, func(ctx context.Context) error { return nil })
	_ = r.Register("job", s)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := r.RunAll(ctx)
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context error, got %v", err)
	}
}
