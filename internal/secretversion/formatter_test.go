package secretversion_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/secretversion"
)

func makeEntry(path string, version int) secretversion.Entry {
	now := time.Now()
	return secretversion.Entry{
		Path:      path,
		Version:   version,
		CreatedAt: now.Add(-24 * time.Hour),
		UpdatedAt: now,
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretversion.FormatTable([]secretversion.Entry{makeEntry("secret/db", 2)})
	for _, h := range []string{"PATH", "VERSION", "CREATED", "UPDATED"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestFormatTable_ContainsEntry(t *testing.T) {
	out := secretversion.FormatTable([]secretversion.Entry{makeEntry("secret/db", 5)})
	if !strings.Contains(out, "secret/db") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "5") {
		t.Error("expected version in output")
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := secretversion.FormatTable(nil)
	if !strings.Contains(out, "no secret versions tracked") {
		t.Error("expected empty message")
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	entries := []secretversion.Entry{
		makeEntry("secret/a", 3),
		makeEntry("secret/b", 7),
	}
	out := secretversion.FormatSummary(entries)
	if !strings.Contains(out, "2") {
		t.Error("expected count 2 in summary")
	}
	if !strings.Contains(out, "7") {
		t.Error("expected highest version 7 in summary")
	}
}

func TestFormatSummary_Empty(t *testing.T) {
	out := secretversion.FormatSummary(nil)
	if !strings.Contains(out, "0") {
		t.Error("expected zero count in summary")
	}
}
