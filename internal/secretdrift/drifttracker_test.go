package secretdrift_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretdrift"
)

var fixedNow = func() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTracker(t *testing.T) *secretdrift.Tracker {
	t.Helper()
	tr, err := secretdrift.New(fixedNow)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestRecord_NoDrift_FirstCall(t *testing.T) {
	tr := newTracker(t)
	if err := tr.Record("secret/foo", "abc123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := len(tr.Drifts()); got != 0 {
		t.Errorf("expected 0 drifts, got %d", got)
	}
}

func TestRecord_DriftDetected_OnHashChange(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Record("secret/foo", "abc123")
	_ = tr.Record("secret/foo", "xyz789")

	drifts := tr.Drifts()
	if len(drifts) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(drifts))
	}
	if drifts[0].PreviousHash != "abc123" {
		t.Errorf("expected prev abc123, got %s", drifts[0].PreviousHash)
	}
	if drifts[0].CurrentHash != "xyz789" {
		t.Errorf("expected curr xyz789, got %s", drifts[0].CurrentHash)
	}
}

func TestRecord_NoDrift_SameHash(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Record("secret/bar", "same")
	_ = tr.Record("secret/bar", "same")
	if got := len(tr.Drifts()); got != 0 {
		t.Errorf("expected 0 drifts, got %d", got)
	}
}

func TestRecord_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	if err := tr.Record("", "hash"); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRecord_EmptyHash_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	if err := tr.Record("secret/x", ""); err == nil {
		t.Error("expected error for empty hash")
	}
}

func TestReset_ClearsDrifts(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Record("secret/foo", "a")
	_ = tr.Record("secret/foo", "b")
	tr.Reset()
	if got := len(tr.Drifts()); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Record("secret/foo", "aaa")
	_ = tr.Record("secret/foo", "bbb")
	out := secretdrift.FormatTable(tr.Drifts())
	for _, h := range []string{"PATH", "PREV HASH", "CURR HASH", "DETECTED AT"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_EmptyDrifts(t *testing.T) {
	out := secretdrift.FormatTable(nil)
	if !strings.Contains(out, "No drift") {
		t.Errorf("expected 'No drift' message, got: %s", out)
	}
}

func TestFormatSummary_Count(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Record("secret/a", "1")
	_ = tr.Record("secret/a", "2")
	out := secretdrift.FormatSummary(tr.Drifts())
	if !strings.Contains(out, "1 change") {
		t.Errorf("expected count in summary, got: %s", out)
	}
}
