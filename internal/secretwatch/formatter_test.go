package secretwatch_test

import (
	"strings"
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/secretwatch"
)

func makeEvents() []secretwatch.Event {
	return []secretwatch.Event{
		{Path: "secret/db", Kind: "added", Detail: ""},
		{Path: "secret/api", Kind: "removed", Detail: ""},
		{Path: "secret/token", Kind: "modified", Detail: "version bump"},
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretwatch.FormatTable(makeEvents())
	for _, h := range []string{"PATH", "CHANGE", "DETAIL"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_ChangeLabels(t *testing.T) {
	out := secretwatch.FormatTable(makeEvents())
	for _, label := range []string{"ADDED", "REMOVED", "MODIFIED"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestFormatTable_EmptyEvents(t *testing.T) {
	out := secretwatch.FormatTable(nil)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected 'No changes' message, got: %q", out)
	}
}

func TestFormatTable_ContainsPaths(t *testing.T) {
	out := secretwatch.FormatTable(makeEvents())
	for _, path := range []string{"secret/db", "secret/api", "secret/token"} {
		if !strings.Contains(out, path) {
			t.Errorf("expected path %q in output", path)
		}
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := secretwatch.FormatSummary(makeEvents())
	for _, s := range []string{"added=1", "removed=1", "modified=1"} {
		if !strings.Contains(out, s) {
			t.Errorf("expected %q in summary, got: %q", s, out)
		}
	}
}
