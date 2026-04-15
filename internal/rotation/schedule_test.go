package rotation

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestEvaluate_NotOverdue(t *testing.T) {
	now := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	eval := NewEvaluator(fixedNow(now))

	s := Schedule{
		Path:        "secret/db",
		Interval:    7 * 24 * time.Hour,
		LastRotated: now.Add(-3 * 24 * time.Hour),
	}
	st, err := eval.Evaluate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if st.Overdue {
		t.Error("expected not overdue")
	}
	if st.DueIn <= 0 {
		t.Errorf("expected positive DueIn, got %v", st.DueIn)
	}
}

func TestEvaluate_Overdue(t *testing.T) {
	now := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	eval := NewEvaluator(fixedNow(now))

	s := Schedule{
		Path:        "secret/api",
		Interval:    24 * time.Hour,
		LastRotated: now.Add(-48 * time.Hour),
	}
	st, err := eval.Evaluate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !st.Overdue {
		t.Error("expected overdue")
	}
	if st.DueIn >= 0 {
		t.Errorf("expected negative DueIn, got %v", st.DueIn)
	}
}

func TestEvaluate_InvalidInterval(t *testing.T) {
	eval := NewEvaluator(nil)
	s := Schedule{Path: "secret/bad", Interval: 0, LastRotated: time.Now()}
	_, err := eval.Evaluate(s)
	if err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestEvaluateAll_PartialErrors(t *testing.T) {
	now := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	eval := NewEvaluator(fixedNow(now))

	schedules := []Schedule{
		{Path: "secret/ok", Interval: 24 * time.Hour, LastRotated: now.Add(-1 * time.Hour)},
		{Path: "secret/bad", Interval: 0, LastRotated: now},
	}
	statuses, errs := eval.EvaluateAll(schedules)
	if len(statuses) != 1 {
		t.Errorf("expected 1 status, got %d", len(statuses))
	}
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %d", len(errs))
	}
}
