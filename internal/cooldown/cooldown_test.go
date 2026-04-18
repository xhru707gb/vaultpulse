package cooldown_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/cooldown"
)

var (
	epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func newTracker(t *testing.T, window time.Duration) *cooldown.Tracker {
	t.Helper()
	tr, err := cooldown.New(window)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	tr.SetNow(fixedNow(epoch))
	return tr
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := cooldown.New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestIsCoolingDown_FirstCall_ReturnsFalse(t *testing.T) {
	tr := newTracker(t, time.Minute)
	if tr.IsCoolingDown("secret/a") {
		t.Fatal("expected false on first call")
	}
}

func TestIsCoolingDown_AfterRecord_ReturnsTrue(t *testing.T) {
	tr := newTracker(t, time.Minute)
	tr.Record("secret/a")
	if !tr.IsCoolingDown("secret/a") {
		t.Fatal("expected true immediately after Record")
	}
}

func TestIsCoolingDown_AfterWindowExpires_ReturnsFalse(t *testing.T) {
	tr := newTracker(t, time.Minute)
	tr.Record("secret/a")
	tr.SetNow(fixedNow(epoch.Add(2 * time.Minute)))
	if tr.IsCoolingDown("secret/a") {
		t.Fatal("expected false after window expires")
	}
}

func TestReset_ClearsCooldown(t *testing.T) {
	tr := newTracker(t, time.Minute)
	tr.Record("secret/a")
	tr.Reset("secret/a")
	if tr.IsCoolingDown("secret/a") {
		t.Fatal("expected false after Reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	tr := newTracker(t, time.Minute)
	tr.Record("secret/a")
	tr.Record("secret/b")
	if got := tr.Len(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}
