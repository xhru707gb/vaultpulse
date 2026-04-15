package rotation

import (
	"strings"
	"testing"
	"time"
)

func makeStatus(path string, dueIn time.Duration) Status {
	now := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	last := now.Add(-7 * 24 * time.Hour)
	return Status{
		Path:        path,
		DueIn:       dueIn,
		Overdue:     dueIn < 0,
		LastRotated: last,
		NextDue:     last.Add(7 * 24 * time.Hour),
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	var sb strings.Builder
	FormatTable(&sb, []Status{})
	out := sb.String()
	for _, h := range []string{"PATH", "STATUS", "LAST ROTATED", "NEXT DUE", "DUE IN"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestFormatTable_RotationLabels(t *testing.T) {
	statuses := []Status{
		makeStatus("secret/ok", 48*time.Hour),
		makeStatus("secret/soon", 12*time.Hour),
		makeStatus("secret/overdue", -2*time.Hour),
	}
	var sb strings.Builder
	FormatTable(&sb, statuses)
	out := sb.String()

	if !strings.Contains(out, labelOK) {
		t.Errorf("expected %q label", labelOK)
	}
	if !strings.Contains(out, labelDueSoon) {
		t.Errorf("expected %q label", labelDueSoon)
	}
	if !strings.Contains(out, labelOverdue) {
		t.Errorf("expected %q label", labelOverdue)
	}
}

func TestFormatTable_ContainsPaths(t *testing.T) {
	paths := []string{"secret/alpha", "secret/beta", "secret/gamma"}
	var statuses []Status
	for _, p := range paths {
		statuses = append(statuses, makeStatus(p, 24*time.Hour))
	}
	var sb strings.Builder
	FormatTable(&sb, statuses)
	out := sb.String()
	for _, p := range paths {
		if !strings.Contains(out, p) {
			t.Errorf("expected path %q in table output", p)
		}
	}
}

func TestFormatDueIn_Positive(t *testing.T) {
	d := 2*time.Hour + 30*time.Minute
	out := formatDueIn(d)
	if !strings.Contains(out, "2h30m") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDueIn_Negative(t *testing.T) {
	d := -3 * time.Hour
	out := formatDueIn(d)
	if !strings.Contains(out, "ago") {
		t.Errorf("expected 'ago' in output, got %q", out)
	}
}
