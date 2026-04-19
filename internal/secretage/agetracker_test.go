package secretage

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestTracker() *Tracker {
	t := New()
	t.now = func() time.Time { return fixedNow }
	return t
}

func TestRegister_And_Evaluate(t *testing.T) {
	tr := newTestTracker()
	created := fixedNow.Add(-10 * 24 * time.Hour)
	_ = tr.Register("secret/db", created, 30*24*time.Hour)
	s, ok := tr.Evaluate("secret/db")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if s.Overdue {
		t.Error("expected not overdue")
	}
}

func TestEvaluate_Overdue(t *testing.T) {
	tr := newTestTracker()
	created := fixedNow.Add(-40 * 24 * time.Hour)
	_ = tr.Register("secret/old", created, 30*24*time.Hour)
	s, _ := tr.Evaluate("secret/old")
	if !s.Overdue {
		t.Error("expected overdue")
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	err := tr.Register("", fixedNow, time.Hour)
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRegister_InvalidMaxAge_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	err := tr.Register("secret/x", fixedNow, 0)
	if err == nil {
		t.Error("expected error for zero maxAge")
	}
}

func TestEvaluateAll_ReturnsAll(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/a", fixedNow.Add(-5*24*time.Hour), 30*24*time.Hour)
	_ = tr.Register("secret/b", fixedNow.Add(-35*24*time.Hour), 30*24*time.Hour)
	statuses := tr.EvaluateAll()
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}
	overdue := 0
	for _, s := range statuses {
		if s.Overdue {
			overdue++
		}
	}
	if overdue != 1 {
		t.Errorf("expected 1 overdue, got %d", overdue)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := FormatTable(nil)
	if out != "no secrets tracked\n" {
		t.Errorf("unexpected empty output: %q", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	statuses := []Status{
		{Path: "a", Overdue: false},
		{Path: "b", Overdue: true},
	}
	s := FormatSummary(statuses)
	if s != "total=2 overdue=1\n" {
		t.Errorf("unexpected summary: %q", s)
	}
}
