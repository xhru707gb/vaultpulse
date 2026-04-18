package jitter_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/jitter"
)

func deterministicJitter(t *testing.T, factor float64) *jitter.Jitter {
	t.Helper()
	j, err := jitter.New(jitter.Config{
		Factor: factor,
		Source: rand.NewSource(42),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return j
}

func TestNew_InvalidFactor(t *testing.T) {
	for _, f := range []float64{-0.1, 1.1, -1} {
		_, err := jitter.New(jitter.Config{Factor: f})
		if err == nil {
			t.Fatalf("expected error for factor %v", f)
		}
	}
}

func TestNew_ValidFactor(t *testing.T) {
	_, err := jitter.New(jitter.Config{Factor: 0.5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApply_ZeroFactor_ReturnsBase(t *testing.T) {
	j := deterministicJitter(t, 0)
	base := 10 * time.Second
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_AddsPositiveDelta(t *testing.T) {
	j := deterministicJitter(t, 0.5)
	base := 10 * time.Second
	got := j.Apply(base)
	if got < base {
		t.Fatalf("Apply should not return less than base: got %v", got)
	}
	if got > base+5*time.Second {
		t.Fatalf("Apply exceeded base+factor*base: got %v", got)
	}
}

func TestApply_NegativeBase_ReturnsBase(t *testing.T) {
	j := deterministicJitter(t, 0.5)
	base := -1 * time.Second
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApplyRange_WithinBounds(t *testing.T) {
	j := deterministicJitter(t, 0.2)
	base := 10 * time.Second
	for i := 0; i < 50; i++ {
		got := j.ApplyRange(base)
		low := time.Duration(float64(base) * 0.8)
		high := time.Duration(float64(base) * 1.2)
		if got < low || got > high {
			t.Fatalf("ApplyRange out of bounds: got %v, want [%v, %v]", got, low, high)
		}
	}
}
