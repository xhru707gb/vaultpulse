package secretaudit

import (
	"testing"
	"time"
)

func fixedNow() func() time.Time {
	t := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func newTestAuditor() *Auditor {
	a := New()
	a.now = fixedNow()
	return a
}

func TestRecord_And_ForPath(t *testing.T) {
	a := newTestAuditor()
	if err := a.Record("secret/db", EventRead, "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := a.ForPath("secret/db")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %s", events[0].Actor)
	}
	if events[0].Kind != EventRead {
		t.Errorf("expected kind read, got %s", events[0].Kind)
	}
}

func TestRecord_EmptyPath_ReturnsError(t *testing.T) {
	a := newTestAuditor()
	if err := a.Record("", EventWrite, "bob"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRecord_EmptyActor_ReturnsError(t *testing.T) {
	a := newTestAuditor()
	if err := a.Record("secret/x", EventWrite, ""); err == nil {
		t.Fatal("expected error for empty actor")
	}
}

func TestForPath_NoMatch_ReturnsNil(t *testing.T) {
	a := newTestAuditor()
	_ = a.Record("secret/db", EventRead, "alice")
	if got := a.ForPath("secret/other"); len(got) != 0 {
		t.Errorf("expected no events, got %d", len(got))
	}
}

func TestAll_ReturnsAllEvents(t *testing.T) {
	a := newTestAuditor()
	_ = a.Record("secret/a", EventRead, "alice")
	_ = a.Record("secret/b", EventWrite, "bob")
	if got := a.All(); len(got) != 2 {
		t.Errorf("expected 2 events, got %d", len(got))
	}
}

func TestReset_ClearsEvents(t *testing.T) {
	a := newTestAuditor()
	_ = a.Record("secret/a", EventRotate, "carol")
	a.Reset()
	if got := a.All(); len(got) != 0 {
		t.Errorf("expected 0 events after reset, got %d", len(got))
	}
}

func TestRecord_TimestampSet(t *testing.T) {
	a := newTestAuditor()
	expected := fixedNow()()
	_ = a.Record("secret/ts", EventDelete, "dave")
	events := a.All()
	if !events[0].Timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, events[0].Timestamp)
	}
}
