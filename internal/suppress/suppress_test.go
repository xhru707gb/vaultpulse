package suppress

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestSuppressor(t *testing.T, window time.Duration) *Suppressor {
	t.Helper()
	s, err := New(window)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	s.now = func() time.Time { return fixedNow }
	return s
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestIsSuppressed_FirstCall_ReturnsFalse(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	if s.IsSuppressed("secret/a") {
		t.Fatal("expected false on first call")
	}
}

func TestIsSuppressed_AfterRecord_ReturnsTrue(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	s.Record("secret/a")
	if !s.IsSuppressed("secret/a") {
		t.Fatal("expected suppressed after Record")
	}
}

func TestIsSuppressed_AfterWindowExpires_ReturnsFalse(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	s.Record("secret/a")
	s.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }
	if s.IsSuppressed("secret/a") {
		t.Fatal("expected false after window expires")
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	s.Record("secret/b")
	s.Record("secret/b")
	all := s.All()
	if all["secret/b"].Count != 2 {
		t.Fatalf("expected count 2, got %d", all["secret/b"].Count)
	}
}

func TestReset_ClearsSuppression(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	s.Record("secret/c")
	s.Reset("secret/c")
	if s.IsSuppressed("secret/c") {
		t.Fatal("expected false after Reset")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := newTestSuppressor(t, time.Minute)
	s.Record("secret/d")
	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	delete(all, "secret/d")
	if len(s.All()) != 1 {
		t.Fatal("All() should return a copy, not the internal map")
	}
}
