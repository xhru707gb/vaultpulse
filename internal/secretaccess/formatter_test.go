package secretaccess

import (
	"strings"
	"testing"
	"time"
)

func makeAccessEntries() []AccessEntry {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return []AccessEntry{
		{Path: "secret/app/db", Count: 5, LastAccess: now, FirstAccess: now.Add(-24 * time.Hour)},
		{Path: "secret/app/api", Count: 12, LastAccess: now.Add(-1 * time.Hour), FirstAccess: now.Add(-48 * time.Hour)},
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := FormatTable(makeAccessEntries())
	for _, hdr := range []string{"PATH", "COUNT", "LAST ACCESS", "FIRST ACCESS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_ContainsEntries(t *testing.T) {
	out := FormatTable(makeAccessEntries())
	if !strings.Contains(out, "secret/app/db") {
		t.Error("expected path secret/app/db in output")
	}
	if !strings.Contains(out, "12") {
		t.Error("expected count 12 in output")
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := FormatTable(nil)
	if !strings.Contains(out, "no access records") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := FormatSummary(makeAccessEntries())
	if !strings.Contains(out, "total paths: 2") {
		t.Errorf("expected total paths: 2, got: %s", out)
	}
	if !strings.Contains(out, "total accesses: 17") {
		t.Errorf("expected total accesses: 17, got: %s", out)
	}
}

func TestFormatSummary_Empty(t *testing.T) {
	out := FormatSummary(nil)
	if !strings.Contains(out, "total paths: 0") {
		t.Errorf("expected zero summary, got: %s", out)
	}
}

func TestFormatTime_Zero(t *testing.T) {
	result := formatTime(time.Time{})
	if result != "never" {
		t.Errorf("expected 'never' for zero time, got %q", result)
	}
}
