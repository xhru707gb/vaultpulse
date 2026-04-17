package sampling_test

import (
	"testing"

	"github.com/your-org/vaultpulse/internal/sampling"
)

func TestNew_InvalidRate(t *testing.T) {
	for _, r := range []float64{-0.1, 1.1, -1.0} {
		_, err := sampling.New(sampling.Config{Rate: r})
		if err == nil {
			t.Fatalf("expected error for rate %v", r)
		}
	}
}

func TestNew_ValidRate(t *testing.T) {
	s, err := sampling.New(sampling.Config{Rate: 0.5, Seed: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Rate() != 0.5 {
		t.Fatalf("expected rate 0.5, got %v", s.Rate())
	}
}

func TestSample_RateZero_AlwaysFalse(t *testing.T) {
	s, _ := sampling.New(sampling.Config{Rate: 0.0, Seed: 1})
	for i := 0; i < 100; i++ {
		if s.Sample() {
			t.Fatal("expected Sample() == false for rate 0")
		}
	}
}

func TestSample_RateOne_AlwaysTrue(t *testing.T) {
	s, _ := sampling.New(sampling.Config{Rate: 1.0, Seed: 1})
	for i := 0; i < 100; i++ {
		if !s.Sample() {
			t.Fatal("expected Sample() == true for rate 1")
		}
	}
}

func TestSample_PartialRate_Approximate(t *testing.T) {
	s, _ := sampling.New(sampling.Config{Rate: 0.5, Seed: 99})
	hits := 0
	const n = 10_000
	for i := 0; i < n; i++ {
		if s.Sample() {
			hits++
		}
	}
	ratio := float64(hits) / float64(n)
	if ratio < 0.40 || ratio > 0.60 {
		t.Fatalf("expected ~50%% hits, got %.2f%%", ratio*100)
	}
}
