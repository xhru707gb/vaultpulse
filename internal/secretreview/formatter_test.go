package secretreview_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretreview"
)

func makeEntries() []*secretreview.Entry {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	return []*secretreview.Entry{
		{
			Path:       "secret/alpha",
			Reviewer:   "alice",
			Status:     secretreview.StatusApproved,
			NextReview: base.Add(24 * time.Hour),
			Interval:   48 * time.Hour,
		},
		{
			Path:       "secret/beta",
			Reviewer:   "bob",
			Status:     secretreview.StatusOverdue,
			NextReview: base.Add(-12 * time.Hour),
			Interval:   24 * time.Hour,
		},
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretreview.FormatTable(makeEntries())
	for _, h := range []string{"PATH", "REVIEWER", "STATUS", "NEXT REVIEW", "INTERVAL"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q in output", h)
		}
	}
}

func TestFormatTable_StatusLabels(t *testing.T) {
	out := secretreview.FormatTable(makeEntries())
	if !strings.Contains(out, "approved") {
		t.Error("expected 'approved' label")
	}
	if !strings.Contains(out, "overdue") {
		t.Error("expected 'overdue' label")
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := secretreview.FormatTable(nil)
	if !strings.Contains(out, "no review entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := secretreview.FormatSummary(makeEntries())
	if !strings.Contains(out, "total=2") {
		t.Errorf("expected total=2 in summary, got: %s", out)
	}
	if !strings.Contains(out, "overdue=1") {
		t.Errorf("expected overdue=1 in summary, got: %s", out)
	}
	if !strings.Contains(out, "approved=1") {
		t.Errorf("expected approved=1 in summary, got: %s", out)
	}
}
