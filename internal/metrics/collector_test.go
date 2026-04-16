package metrics

import (
	"errors"
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func TestRecord_StoresSnapshot(t *testing.T) {
	c := NewCollector()
	c.now = func() time.Time { return fixedNow }

	c.Record("secret/db", "ok", 42*time.Millisecond, nil)

	s, ok := c.Get("secret/db")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if s.Status != "ok" {
		t.Errorf("status: got %q, want %q", s.Status, "ok")
	}
	if s.Duration != 42*time.Millisecond {
		t.Errorf("duration: got %v, want 42ms", s.Duration)
	}
	if !s.CheckedAt.Equal(fixedNow) {
		t.Errorf("checkedAt: got %v, want %v", s.CheckedAt, fixedNow)
	}
	if s.Error != "" {
		t.Errorf("error: expected empty, got %q", s.Error)
	}
}

func TestRecord_OverwritesPrevious(t *testing.T) {
	c := NewCollector()
	c.Record("secret/db", "warning", 10*time.Millisecond, nil)
	c.Record("secret/db", "expired", 20*time.Millisecond, errors.New("ttl elapsed"))

	s, _ := c.Get("secret/db")
	if s.Status != "expired" {
		t.Errorf("expected overwritten status %q, got %q", "expired", s.Status)
	}
	if s.Error != "ttl elapsed" {
		t.Errorf("expected error %q, got %q", "ttl elapsed", s.Error)
	}
}

func TestReset_ClearsSnapshot(t *testing.T) {
	c := NewCollector()
	c.Record("secret/db", "ok", 0, nil)
	c.Reset()

	if _, ok := c.Get("secret/db"); ok {
		t.Error("expected snapshot to be cleared after Reset")
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := NewCollector()
	c.Record("secret/api", "ok", 5*time.Millisecond, nil)

	s1, _ := c.Get("secret/api")
	s1.Status = "mutated"

	s2, _ := c.Get("secret/api")
	if s2.Status == "mutated" {
		t.Error("Get should return a value copy, not a pointer to internal state")
	}
}
