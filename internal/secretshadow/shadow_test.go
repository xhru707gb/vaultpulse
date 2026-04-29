package secretshadow

import (
	"strings"
	"testing"
	"time"
)

func newTestTracker() *Tracker {
	t := New()
	fixed := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	t.nowFn = func() time.Time { return fixed }
	return t
}

func TestCapture_And_CheckNoDivergence(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Capture("secret/foo", "value1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := tr.Check("secret/foo", "value1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Diverged {
		t.Error("expected no divergence for same value")
	}
}

func TestCheck_Diverged(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Capture("secret/bar", "original")
	e, err := tr.Check("secret/bar", "modified")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !e.Diverged {
		t.Error("expected divergence for changed value")
	}
}

func TestCapture_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Capture("", "v"); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestCheck_NotFound_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	_, err := tr.Check("secret/unknown", "v")
	if err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Capture("secret/baz", "val")
	if err := tr.Remove("secret/baz"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tr.All()) != 0 {
		t.Error("expected empty tracker after remove")
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Remove("secret/nope"); err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Capture("secret/x", "abc")
	entry, _ := tr.Check("secret/x", "abc")
	out := FormatTable([]Entry{entry})
	for _, h := range []string{"PATH", "HASH", "CAPTURED AT", "STATUS"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	entries := []Entry{
		{Path: "a", Diverged: false},
		{Path: "b", Diverged: true},
		{Path: "c", Diverged: true},
	}
	out := FormatSummary(entries)
	if !strings.Contains(out, "3 total") {
		t.Errorf("expected total count, got: %s", out)
	}
	if !strings.Contains(out, "2 diverged") {
		t.Errorf("expected diverged count, got: %s", out)
	}
}
