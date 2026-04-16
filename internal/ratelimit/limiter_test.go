package ratelimit

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New(Config{RequestsPerSecond: 0, Burst: 10})
	if err == nil {
		t.Fatal("expected error for zero RequestsPerSecond")
	}

	_, err = New(Config{RequestsPerSecond: 10, Burst: 0})
	if err == nil {
		t.Fatal("expected error for zero Burst")
	}
}

func TestAllow_WithinBurst(t *testing.T) {
	l, err := New(Config{RequestsPerSecond: 10, Burst: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := time.Now()
	l.now = fixedNow(base)
	l.lastTick = base

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	// Burst exhausted — next call should be denied
	if l.Allow() {
		t.Fatal("expected Allow()=false after burst exhausted")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l, err := New(Config{RequestsPerSecond: 10, Burst: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := time.Now()
	l.now = fixedNow(base)
	l.lastTick = base

	// Consume the single token
	if !l.Allow() {
		t.Fatal("expected first Allow()=true")
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after token consumed")
	}

	// Advance clock by 200ms — should refill 2 tokens (rate=10/s), capped at burst=1
	l.now = fixedNow(base.Add(200 * time.Millisecond))
	if !l.Allow() {
		t.Fatal("expected Allow()=true after refill")
	}
}

func TestAllow_TokensCappedAtBurst(t *testing.T) {
	l, err := New(Config{RequestsPerSecond: 100, Burst: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	base := time.Now()
	l.now = fixedNow(base)
	l.lastTick = base

	// Advance a full second — would add 100 tokens but capped at 2
	l.now = fixedNow(base.Add(time.Second))

	allowed := 0
	for i := 0; i < 5; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 2 {
		t.Fatalf("expected 2 allowed after refill cap, got %d", allowed)
	}
}
