package metrics

import (
	"strings"
	"testing"
	"time"
)

func makeSnapshot() Snapshot {
	return Snapshot{
		CollectedAt:     time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		LastCheckDur:    250 * time.Millisecond,
		TotalSecrets:    10,
		Expired:         2,
		Warning:         3,
		Healthy:         5,
		OverdueRotation: 1,
		PolicyViolation: 2,
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := FormatTable(makeSnapshot())
	for _, want := range []string{"Metric", "Value", "VaultPulse Metrics Snapshot"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected header %q in output", want)
		}
	}
}

func TestFormatTable_ContainsValues(t *testing.T) {
	out := FormatTable(makeSnapshot())
	for _, want := range []string{"10", "2", "3", "5", "1", "250ms"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected value %q in output", want)
		}
	}
}

func TestFormatTable_ZeroTime(t *testing.T) {
	s := Snapshot{}
	out := FormatTable(s)
	if !strings.Contains(out, "n/a") {
		t.Error("expected 'n/a' for zero CollectedAt")
	}
}

func TestFormatTable_ZeroDuration(t *testing.T) {
	s := Snapshot{CollectedAt: time.Now()}
	out := FormatTable(s)
	if !strings.Contains(out, "n/a") {
		t.Error("expected 'n/a' for zero LastCheckDur")
	}
}
