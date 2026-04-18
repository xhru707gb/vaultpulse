package tagindex_test

import (
	"strings"
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/tagindex"
)

func buildIndex(t *testing.T) *tagindex.Index {
	t.Helper()
	idx := tagindex.New()
	_ = idx.Add("env:prod", "secret/db")
	_ = idx.Add("env:prod", "secret/api")
	_ = idx.Add("team:sre", "secret/infra")
	return idx
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := tagindex.FormatTable(buildIndex(t))
	for _, h := range []string{"TAG", "PATHS", "MEMBERS"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_ContainsTags(t *testing.T) {
	out := tagindex.FormatTable(buildIndex(t))
	if !strings.Contains(out, "env:prod") {
		t.Error("expected env:prod in output")
	}
	if !strings.Contains(out, "team:sre") {
		t.Error("expected team:sre in output")
	}
}

func TestFormatTable_EmptyIndex(t *testing.T) {
	idx := tagindex.New()
	out := tagindex.FormatTable(idx)
	if !strings.Contains(out, "no tags") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatTable_TruncatesLongPaths(t *testing.T) {
	idx := tagindex.New()
	for i := 0; i < 5; i++ {
		_ = idx.Add("big", strings.Repeat("x", 10)+string(rune('a'+i)))
	}
	out := tagindex.FormatTable(idx)
	if !strings.Contains(out, "+2 more") {
		t.Errorf("expected truncation indicator, got: %s", out)
	}
}
