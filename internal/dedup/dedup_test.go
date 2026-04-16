package dedup_test

import (
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/dedup"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestNew_InvalidWindow(t *testing.T) {
	_, err := dedup.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestIsDuplicate_FirstCall_ReturnsFalse(t *testing.T) {
	d, _ := dedup.New(10 * time.Minute)
	d.SetNow(fixedNow(epoch))
	if d.IsDuplicate("secret/a", "expired") {
		t.Fatal("first call should not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindow_ReturnsTrue(t *testing.T) {
	d, _ := dedup.New(10 * time.Minute)
	d.SetNow(fixedNow(epoch))
	d.IsDuplicate("secret/a", "expired")
	if !d.IsDuplicate("secret/a", "expired") {
		t.Fatal("second call within window should be duplicate")
	}
}

func TestIsDuplicate_AfterWindowExpires_ReturnsFalse(t *testing.T) {
	d, _ := dedup.New(10 * time.Minute)
	d.SetNow(fixedNow(epoch))
	d.IsDuplicate("secret/a", "expired")
	d.SetNow(fixedNow(epoch.Add(11 * time.Minute)))
	if d.IsDuplicate("secret/a", "expired") {
		t.Fatal("call after window should not be duplicate")
	}
}

func TestIsDuplicate_DifferentEvent_ReturnsFalse(t *testing.T) {
	d, _ := dedup.New(10 * time.Minute)
	d.SetNow(fixedNow(epoch))
	d.IsDuplicate("secret/a", "expired")
	if d.IsDuplicate("secret/a", "warning") {
		t.Fatal("different event should not be duplicate")
	}
}

func TestFlush_RemovesExpiredEntries(t *testing.T) {
	d, _ := dedup.New(5 * time.Minute)
	d.SetNow(fixedNow(epoch))
	d.IsDuplicate("secret/a", "expired")
	d.IsDuplicate("secret/b", "warning")
	d.SetNow(fixedNow(epoch.Add(6 * time.Minute)))
	n := d.Flush()
	if n != 2 {
		t.Fatalf("expected 2 flushed, got %d", n)
	}
	if d.Len() != 0 {
		t.Fatalf("expected 0 remaining, got %d", d.Len())
	}
}
