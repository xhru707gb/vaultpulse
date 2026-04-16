package throttle_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/throttle"
)

func TestNew_InvalidConcurrency(t *testing.T) {
	_, err := throttle.New(0)
	if err == nil {
		t.Fatal("expected error for concurrency=0")
	}
}

func TestNew_ValidConcurrency(t *testing.T) {
	th, err := throttle.New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.Capacity() != 5 {
		t.Fatalf("expected capacity 5, got %d", th.Capacity())
	}
}

func TestAcquireRelease(t *testing.T) {
	th, _ := throttle.New(2)
	ctx := context.Background()

	if err := th.Acquire(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if th.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", th.InFlight())
	}
	th.Release()
	if th.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after release")
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	th, _ := throttle.New(1)
	ctx := context.Background()

	// fill the slot
	_ = th.Acquire(ctx)

	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := th.Acquire(ctx2)
	if err != throttle.ErrThrottled {
		t.Fatalf("expected ErrThrottled, got %v", err)
	}
}

func TestDo_LimitsConcurrency(t *testing.T) {
	const limit = 3
	th, _ := throttle.New(limit)

	var peak int64
	var current int64
	var wg sync.WaitGroup

	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = th.Do(context.Background(), func() error {
				v := atomic.AddInt64(&current, 1)
				for {
					p := atomic.LoadInt64(&peak)
					if v <= p || atomic.CompareAndSwapInt64(&peak, p, v) {
						break
					}
				}
				time.Sleep(10 * time.Millisecond)
				atomic.AddInt64(&current, -1)
				return nil
			})
		}()
	}
	wg.Wait()

	if peak > int64(limit) {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
