package semaphore_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/semaphore"
)

func TestNew_InvalidSize(t *testing.T) {
	_, err := semaphore.New(0)
	if err == nil {
		t.Fatal("expected error for size 0")
	}
}

func TestNew_ValidSize(t *testing.T) {
	sem, err := semaphore.New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sem.Size() != 5 {
		t.Fatalf("expected size 5, got %d", sem.Size())
	}
}

func TestAcquireRelease(t *testing.T) {
	sem, _ := semaphore.New(2)
	ctx := context.Background()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("acquire 1: %v", err)
	}
	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("acquire 2: %v", err)
	}
	if sem.Acquired() != 2 {
		t.Fatalf("expected 2 acquired, got %d", sem.Acquired())
	}
	sem.Release()
	if sem.Acquired() != 1 {
		t.Fatalf("expected 1 acquired after release, got %d", sem.Acquired())
	}
}

func TestAcquire_ContextCancelled(t *testing.T) {
	sem, _ := semaphore.New(1)
	ctx := context.Background()
	_ = sem.Acquire(ctx) // fill the slot

	cancel, cancelFn := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancelFn()

	err := sem.Acquire(cancel)
	if err == nil {
		t.Fatal("expected context error when semaphore is full")
	}
}

func TestConcurrentAcquire_BoundedBySize(t *testing.T) {
	const size = 3
	const goroutines = 10
	sem, _ := semaphore.New(size)

	var mu sync.Mutex
	peak := 0
	current := 0
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sem.Acquire(context.Background())
			mu.Lock()
			current++
			if current > peak {
				peak = current
			}
			mu.Unlock()
			time.Sleep(5 * time.Millisecond)
			mu.Lock()
			current--
			mu.Unlock()
			sem.Release()
		}()
	}
	wg.Wait()
	if peak > size {
		t.Fatalf("peak concurrency %d exceeded semaphore size %d", peak, size)
	}
}
