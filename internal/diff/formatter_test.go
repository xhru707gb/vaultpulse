package diff

import (
	"strings"
	"testing"
	"time"
)

func makeChange(path string, kind ChangeKind, old, new interface{}) ChangeEntry {
	return ChangeEntry{
		Path:     path,
		Kind:     kind,
		OldValue: old,
		NewValue: new,
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	changes := []ChangeEntry{
		makeChange("secret/foo", KindAdded, nil, "v2"),
	}
	out := FormatTable(changes)
	for _, hdr := range []string{"PATH", "CHANGE", "OLD VALUE", "NEW VALUE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_ChangeLabels(t *testing.T) {
	cases := []struct {
		kind  ChangeKind
		label string
	}{
		{KindAdded, "ADDED"},
		{"},
		{KindModified, "MODIFIED"},
	}
	for _, tc := range cases {
		changes := []ChangeEntry{makeChange("secret/x", tc.kind, nil, nil)}
		out := FormatTable(changes)
		if !strings.Contains(out, tc.label) {
			t.Errorf("expected label %q for kind %v", tc.label, tc.kind)
		}
	}
}

func TestFormatTable_EmptyChanges(t *testing.T) {
	out := FormatTable(nil)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", out)
	}
}

func TestFormatTable_TimeField(t *testing.T) {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	changes := []ChangeEntry{
		makeChange("secret/bar", KindModified, ts, ts.Add(24*time.Hour)),
	}
	out := FormatTable(changes)
	if !strings.Contains(out, "2024-06-01T12:00:00Z") {
		t.Errorf("expected formatted timestamp in got:\n%s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	changes := []ChangeEntry{
		makeChange("a", KindAdded, nil, "v1"),
		makeChange("b", KindAdded, nil, "v2"),
		makeChange("c", KindRemoved, "v1", nil),
		makeChange("d", KindModified, "v1", "v2"),
	}
	summary := FormatSummary(changes)
	if !strings.Contains(summary, "+2 added") {
		t.Errorf("expected '+2 added' in summary: %s", summary)
	}
	if !strings.Contains(summary, "-1 removed") {
		t.Errorf("expected '-1 removed' in summary: %s", summary)
	}
	if !strings.Contains(summary, "~1 modified") {
		t.Errorf("expected '~1 modified' in summary: %s", summary)
	}
}

func TestTruncatePath_Long(t *testing.T) {
	long := strings.Repeat("a", 50)
	result := truncatePath(long, 40)
	if len([]rune(result)) > 41 {
		t.Errorf("truncated path too long: %d chars", len(result))
	}
}
