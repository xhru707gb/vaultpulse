package presign_test

import (
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/presign"
)

func newTracker(now time.Time) *presign.Tracker {
	t := presign.New()
	presign.SetClock(t, func() time.Time { return now })
	return t
}

var fixed = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestRegister_And_Get(t *testing.T) {
	tr := newTracker(fixed)
	err := tr.Register("secret/app", "tok123", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("secret/app")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Token != "tok123" {
		t.Errorf("expected tok123, got %s", e.Token)
	}
	if e.ExpiresAt != fixed.Add(time.Hour) {
		t.Errorf("unexpected expiry: %v", e.ExpiresAt)
	}
}

func TestRegister_InvalidArgs(t *testing.T) {
	tr := presign.New()
	if err := tr.Register("", "tok", time.Hour); err == nil {
		t.Error("expected error for empty path")
	}
	if err := tr.Register("secret/x", "", time.Hour); err == nil {
		t.Error("expected error for empty token")
	}
	if err := tr.Register("secret/x", "tok", 0); err == nil {
		t.Error("expected error for zero ttl")
	}
}

func TestIsExpired_And_TTL(t *testing.T) {
	tr := newTracker(fixed)
	_ = tr.Register("secret/exp", "tok", time.Minute)
	e, _ := tr.Get("secret/exp")

	if e.IsExpired(fixed) {
		t.Error("should not be expired at issue time")
	}
	if e.IsExpired(fixed.Add(2 * time.Minute)) {
	} else {
		t.Error("should be expired after ttl")
	}
	ttl := e.TTL(fixed.Add(30 * time.Second))
	if ttl <= 0 {
		t.Errorf("expected positive TTL, got %v", ttl)
	}
}

func TestRevoke_RemovesEntry(t *testing.T) {
	tr := presign.New()
	_ = tr.Register("secret/r", "tok", time.Hour)
	if err := tr.Revoke("secret/r"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := tr.Get("secret/r"); ok {
		t.Error("expected entry to be removed")
	}
}

func TestRevoke_Unknown_ReturnsError(t *testing.T) {
	tr := presign.New()
	if err := tr.Revoke("secret/missing"); err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestAll_ReturnsEntries(t *testing.T) {
	tr := presign.New()
	_ = tr.Register("a", "t1", time.Hour)
	_ = tr.Register("b", "t2", time.Hour)
	if len(tr.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(tr.All()))
	}
}
