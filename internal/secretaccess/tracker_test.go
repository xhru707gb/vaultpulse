package secretaccess_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretaccess"
)

func newTestTracker() *secretaccess.Tracker {
	return secretaccess.New()
}

func TestRecord_And_Get(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Record("secret/foo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("secret/foo")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.AccessCount != 1 {
		t.Fatalf("expected count 1, got %d", e.AccessCount)
	}
}

func TestRecord_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Record(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRecord_IncrementsCount(t *testing.T) {
	tr := newTestTracker()
	for i := 0; i < 5; i++ {
		_ = tr.Record("secret/bar")
	}
	e, _ := tr.Get("secret/bar")
	if e.AccessCount != 5 {
		t.Fatalf("expected 5, got %d", e.AccessCount)
	}
}

func TestRecord_UpdatesLastAccess(t *testing.T) {
	tr := newTestTracker()
	before := time.Now()
	_ = tr.Record("secret/ts")
	e, _ := tr.Get("secret/ts")
	if e.LastAccess.Before(before) {
		t.Fatal("LastAccess should be >= before")
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := newTestTracker()
	_, ok := tr.Get("secret/missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Record("secret/a")
	_ = tr.Record("secret/b")
	if got := len(tr.All()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Record("secret/x")
	tr.Reset()
	if got := len(tr.All()); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
