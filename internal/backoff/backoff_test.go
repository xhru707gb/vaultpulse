package backoff_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/backoff"
)

func TestNew_InvalidConfig(t *testing.T) {
	cases := []struct {
		name string
		cfg  backoff.Config
	}{
		{"zero initial", backoff.Config{InitialInterval: 0, Multiplier: 2, MaxInterval: time.Second}},
		{"multiplier < 1", backoff.Config{InitialInterval: time.Millisecond, Multiplier: 0.5, MaxInterval: time.Second}},
		{"max < initial", backoff.Config{InitialInterval: time.Second, Multiplier: 2, MaxInterval: time.Millisecond}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := backoff.New(tc.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestNext_Increases(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      2.0,
		MaxInterval:     10 * time.Second,
		Jitter:          false,
	}
	b, err := backoff.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prev := b.Next()
	for i := 0; i < 5; i++ {
		next := b.Next()
		if next <= prev {
			t.Errorf("expected next (%v) > prev (%v)", next, prev)
		}
		prev = next
	}
}

func TestNext_CapsAtMaxInterval(t *testing.T) {
	cfg := backoff.Config{
		InitialInterval: 100 * time.Millisecond,
		Multiplier:      10.0,
		MaxInterval:     500 * time.Millisecond,
		Jitter:          false,
	}
	b, _ := backoff.New(cfg)
	for i := 0; i < 10; i++ {
		d := b.Next()
		if d > cfg.MaxInterval {
			t.Errorf("duration %v exceeds MaxInterval %v", d, cfg.MaxInterval)
		}
	}
}

func TestReset_ResetsAttempt(t *testing.T) {
	cfg := backoff.DefaultConfig()
	cfg.Jitter = false
	b, _ := backoff.New(cfg)
	first := b.Next()
	b.Next()
	b.Next()
	b.Reset()
	if b.Attempt() != 0 {
		t.Errorf("expected attempt 0 after reset, got %d", b.Attempt())
	}
	if b.Next() != first {
		t.Error("expected first duration after reset")
	}
}

func TestDefaultConfig_Valid(t *testing.T) {
	_, err := backoff.New(backoff.DefaultConfig())
	if err != nil {
		t.Fatalf("DefaultConfig should be valid: %v", err)
	}
}
