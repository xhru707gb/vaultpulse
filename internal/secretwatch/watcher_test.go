package secretwatch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/secretwatch"
)

func staticLister(m map[string]int) secretwatch.SecretLister {
	return func(_ context.Context) (map[string]int, error) {
		copy := make(map[string]int, len(m))
		for k, v := range m {
			copy[k] = v
		}
		return copy, nil
	}
}

func TestNew_InvalidInterval(t *testing.T) {
	_, err := secretwatch.New(staticLister(nil), func(_ []secretwatch.Event) {}, 0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_NilLister(t *testing.T) {
	_, err := secretwatch.New(nil, func(_ []secretwatch.Event) {}, time.Second)
	if err == nil {
		t.Fatal("expected error for nil lister")
	}
}

func TestNew_NilHandler(t *testing.T) {
	_, err := secretwatch.New(staticLister(nil), nil, time.Second)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRun_DetectsAddedSecret(t *testing.T) {
	calls := 0
	var mu sync.Mutex
	var got []secretwatch.Event

	lister := staticLister(map[string]int{"secret/a": 1})
	handler := func(events []secretwatch.Event) {
		mu.Lock()
		defer mu.Unlock()
		got = append(got, events...)
		calls++
	}

	w, err := secretwatch.New(lister, handler, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	mu.Lock()
	defer mu.Unlock()
	if calls == 0 {
		t.Fatal("handler was never called")
	}
	if len(got) == 0 {
		t.Fatal("expected at least one event")
	}
	if got[0].Kind != "added" {
		t.Fatalf("expected kind=added, got %q", got[0].Kind)
	}
}

func TestRun_NoEventsWhenUnchanged(t *testing.T) {
	var mu sync.Mutex
	calls := 0

	// lister always returns the same map
	state := map[string]int{"secret/stable": 2}
	lister := staticLister(state)

	var callCount int
	listerWrapped := func(ctx context.Context) (map[string]int, error) {
		mu.Lock()
		callCount++
		mu.Unlock()
		return lister(ctx)
	}

	handler := func(_ []secretwatch.Event) {
		mu.Lock()
		calls++
		mu.Unlock()
	}

	w, err := secretwatch.New(listerWrapped, handler, 20*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck

	mu.Lock()
	defer mu.Unlock()
	// After first tick handler fires (added). Subsequent ticks: no change.
	if calls > 1 {
		t.Fatalf("expected handler called at most once, got %d", calls)
	}
}
