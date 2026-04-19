package secretexpiry

import (
	"strings"
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
	e := Entry{Path: "secret/a", ExpiresAt: fixedNow.Add(48 * time.Hour), WarnBefore: 24 * time.Hour}
	if err := tr.Register(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := tr.Evaluate("secret/a")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if s.Expired || s.Warning {
		t.Errorf("expected OK status, got expired=%v warning=%v", s.Expired, s.Warning)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	tr := newTestTracker()
	e := Entry{Path: "secret/b", ExpiresAt: fixedNow.Add(6 * time.Hour), WarnBefore: 24 * time.Hour}
	_ = tr.Register(e)
	s, _ := tr.Evaluate("secret/b")
	if !s.Warning {
		t.Error("expected warning status")
	}
}

func TestEvaluate_Expired(t *testing.T) {
	tr := newTestTracker()
	e := Entry{Path: "secret/c", ExpiresAt: fixedNow.Add(-1 * time.Hour), WarnBefore: 24 * time.Hour}
	_ = tr.Register(e)
	s, _ := tr.Evaluate("secret/c")
	if !s.Expired {
		t.Error("expected expired status")
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	err := tr.Register(Entry{ExpiresAt: fixedNow.Add(time.Hour)})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register(Entry{Path: "secret/x", ExpiresAt: fixedNow.Add(48 * time.Hour), WarnBefore: time.Hour})
	out := FormatTable(tr.EvaluateAll())
	if !strings.Contains(out, "PATH") || !strings.Contains(out, "STATE") {
		t.Errorf("missing headers in output: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	statuses := []Status{
		{Expired: true},
		{Warning: true},
		{},
	}
	summary := FormatSummary(statuses)
	if !strings.Contains(summary, "expired=1") || !strings.Contains(summary, "warning=1") {
		t.Errorf("unexpected summary: %s", summary)
	}
}
