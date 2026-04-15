package metrics

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_StoresSnapshot(t *testing.T) {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	c := &Collector{nowFn: fixedNow(now)}

	c.Record(Snapshot{TotalSecrets: 5, Expired: 1, Warning: 2})
	got := c.Get()

	if got.TotalSecrets != 5 {
		t.Errorf("expected TotalSecrets=5, got %d", got.TotalSecrets)
	}
	if got.Expired != 1 {
		t.Errorf("expected Expired=1, got %d", got.Expired)
	}
	if !got.CollectedAt.Equal(now) {
		t.Errorf("expected CollectedAt=%v, got %v", now, got.CollectedAt)
	}
}

func TestRecord_OverwritesPrevious(t *testing.T) {
	c := NewCollector()
	c.Record(Snapshot{TotalSecrets: 3})
	c.Record(Snapshot{TotalSecrets: 7})

	if got := c.Get().TotalSecrets; got != 7 {
		t.Errorf("expected 7, got %d", got)
	}
}

func TestReset_ClearsSnapshot(t *testing.T) {
	c := NewCollector()
	c.Record(Snapshot{TotalSecrets: 4, Expired: 2})
	c.Reset()
	got := c.Get()

	if got.TotalSecrets != 0 || got.Expired != 0 {
		t.Errorf("expected zeroed snapshot after Reset, got %+v", got)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	c := NewCollector()
	c.Record(Snapshot{TotalSecrets: 5})
	s := c.Get()
	s.TotalSecrets = 99

	if c.Get().TotalSecrets == 99 {
		t.Error("Get should return a copy, not a reference")
	}
}
