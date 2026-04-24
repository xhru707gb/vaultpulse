package secretstatus_test

import (
	"strings"
	"testing"
	"time"

	"vaultpulse/internal/secretstatus"
)

func makeEntry(path string, level secretstatus.Level, reasons ...string) *secretstatus.Entry {
	return &secretstatus.Entry{
		Path:        path,
		Level:       level,
		Reasons:     reasons,
		EvaluatedAt: time.Now().UTC(),
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretstatus.FormatTable([]*secretstatus.Entry{
		makeEntry("secret/foo", secretstatus.LevelOK),
	})
	for _, hdr := range []string{"PATH", "STATUS", "REASONS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_StatusLabels(t *testing.T) {
	entries := []*secretstatus.Entry{
		makeEntry("secret/ok", secretstatus.LevelOK),
		makeEntry("secret/warn", secretstatus.LevelWarning, "expires soon"),
		makeEntry("secret/crit", secretstatus.LevelCritical, "overdue"),
	}
	out := secretstatus.FormatTable(entries)
	for _, label := range []string{"OK", "WARNING", "CRITICAL"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := secretstatus.FormatTable(nil)
	if !strings.Contains(out, "No secret status") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	entries := []*secretstatus.Entry{
		makeEntry("a", secretstatus.LevelOK),
		makeEntry("b", secretstatus.LevelWarning),
		makeEntry("c", secretstatus.LevelCritical),
		makeEntry("d", secretstatus.LevelCritical),
	}
	out := secretstatus.FormatSummary(entries)
	for _, sub := range []string{"Total: 4", "OK: 1", "Warning: 1", "Critical: 2"} {
		if !strings.Contains(out, sub) {
			t.Errorf("expected %q in summary, got: %s", sub, out)
		}
	}
}
