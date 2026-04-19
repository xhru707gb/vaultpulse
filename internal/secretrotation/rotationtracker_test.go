package secretrotation_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/secretrotation"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestTracker() *secretrotation.Tracker {
	t, _ := secretrotation.New(func() time.Time { return fixedNow })
	return t
}

func TestRegister_And_Get(t *testing.T) {
	tr := newTestTracker()
	last := fixedNow.Add(-24 * time.Hour)
	if err := tr.Register("secret/a", 48*time.Hour, last); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("secret/a")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Path != "secret/a" {
		t.Errorf("expected path secret/a, got %s", e.Path)
	}
	if e.Overdue {
		t.Error("expected not overdue")
	}
}

func TestRegister_Overdue(t *testing.T) {
	tr := newTestTracker()
	last := fixedNow.Add(-72 * time.Hour)
	if err := tr.Register("secret/b", 48*time.Hour, last); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, _ := tr.Get("secret/b")
	if !e.Overdue {
		t.Error("expected overdue")
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	err := tr.Register("", 24*time.Hour, fixedNow)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRegister_InvalidInterval_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	err := tr.Register("secret/c", 0, fixedNow)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestOverdueCount(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/ok", 48*time.Hour, fixedNow.Add(-24*time.Hour))
	_ = tr.Register("secret/late", 12*time.Hour, fixedNow.Add(-24*time.Hour))
	if got := tr.OverdueCount(); got != 1 {
		t.Errorf("expected 1 overdue, got %d", got)
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/x", time.Hour, fixedNow)
	_ = tr.Register("secret/y", time.Hour, fixedNow)
	if len(tr.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(tr.All()))
	}
}
