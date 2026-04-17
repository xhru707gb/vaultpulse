package rollup

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestAggregator(t *testing.T) *Aggregator {
	t.Helper()
	a, err := New(5 * time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.now = func() time.Time { return fixedNow }
	return a
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestAdd_IncreasesLen(t *testing.T) {
	a := newTestAggregator(t)
	a.Add(Event{Path: "secret/a", Level: "ok"})
	a.Add(Event{Path: "secret/b", Level: "warning"})
	if a.Len() != 2 {
		t.Fatalf("expected 2, got %d", a.Len())
	}
}

func TestFlush_ClearsBuffer(t *testing.T) {
	a := newTestAggregator(t)
	a.Add(Event{Path: "secret/a", Level: "expired"})
	s := a.Flush()
	if s.Total != 1 {
		t.Fatalf("expected total 1, got %d", s.Total)
	}
	if a.Len() != 0 {
		t.Fatal("buffer should be empty after flush")
	}
}

func TestFlush_ByLevelCounts(t *testing.T) {
	a := newTestAggregator(t)
	a.Add(Event{Level: "ok"})
	a.Add(Event{Level: "warning"})
	a.Add(Event{Level: "warning"})
	a.Add(Event{Level: "expired"})
	s := a.Flush()
	if s.ByLevel["ok"] != 1 {
		t.Errorf("expected ok=1, got %d", s.ByLevel["ok"])
	}
	if s.ByLevel["warning"] != 2 {
		t.Errorf("expected warning=2, got %d", s.ByLevel["warning"])
	}
	if s.ByLevel["expired"] != 1 {
		t.Errorf("expected expired=1, got %d", s.ByLevel["expired"])
	}
}

func TestFlush_SetsTimestamp(t *testing.T) {
	a := newTestAggregator(t)
	s := a.Flush()
	if !s.FlushedAt.Equal(fixedNow) {
		t.Errorf("expected %v, got %v", fixedNow, s.FlushedAt)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	a := newTestAggregator(t)
	a.Add(Event{Path: "secret/db", Level: "warning", Message: "expiring soon"})
	s := a.Flush()
	out := FormatTable(s)
	for _, want := range []string{"PATH", "LEVEL", "MESSAGE", "secret/db", "warning"} {
		if !contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
