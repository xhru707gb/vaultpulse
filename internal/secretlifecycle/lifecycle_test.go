package secretlifecycle

import (
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestTracker() *Tracker {
	t := New()
	t.now = fixedNow
	return t
}

func baseEntry(path string) Entry {
	now := fixedNow()
	return Entry{
		Path:        path,
		CreatedAt:   now.Add(-48 * time.Hour),
		LastRotated: now.Add(-10 * time.Hour),
		LastAccess:  now.Add(-1 * time.Hour),
		MaxAge:      24 * time.Hour,
		WarnBefore:  2 * time.Hour,
	}
}

func TestRegister_And_Evaluate(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("secret/foo")
	e.ExpiresAt = fixedNow().Add(10 * time.Hour)
	if err := tr.Register(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := tr.Evaluate("secret/foo")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if s.State != StateActive {
		t.Errorf("expected active, got %s", s.State)
	}
}

func TestEvaluate_Expired(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("secret/expired")
	e.ExpiresAt = fixedNow().Add(-1 * time.Hour)
	_ = tr.Register(e)
	s, _ := tr.Evaluate("secret/expired")
	if s.State != StateExpired {
		t.Errorf("expected expired, got %s", s.State)
	}
}

func TestEvaluate_Expiring(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("secret/expiring")
	e.ExpiresAt = fixedNow().Add(1 * time.Hour) // within WarnBefore=2h
	_ = tr.Register(e)
	s, _ := tr.Evaluate("secret/expiring")
	if s.State != StateExpiring {
		t.Errorf("expected expiring, got %s", s.State)
	}
}

func TestEvaluate_Stale(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("secret/stale")
	e.LastRotated = fixedNow().Add(-30 * time.Hour) // exceeds MaxAge=24h
	_ = tr.Register(e)
	s, _ := tr.Evaluate("secret/stale")
	if s.State != StateStale {
		t.Errorf("expected stale, got %s", s.State)
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("")
	if err := tr.Register(e); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRegister_InvalidMaxAge_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	e := baseEntry("secret/bad")
	e.MaxAge = 0
	if err := tr.Register(e); err == nil {
		t.Error("expected error for zero max age")
	}
}

func TestEvaluateAll_ReturnsAllStatuses(t *testing.T) {
	tr := newTestTracker()
	for _, p := range []string{"secret/a", "secret/b", "secret/c"} {
		e := baseEntry(p)
		e.ExpiresAt = fixedNow().Add(5 * time.Hour)
		_ = tr.Register(e)
	}
	all := tr.EvaluateAll()
	if len(all) != 3 {
		t.Errorf("expected 3 statuses, got %d", len(all))
	}
}
