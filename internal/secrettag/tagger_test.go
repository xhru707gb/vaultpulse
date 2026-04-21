package secrettag_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secrettag"
)

func newTagger() *secrettag.Tagger { return secrettag.New() }

func TestAdd_And_Tags(t *testing.T) {
	tr := newTagger()
	if err := tr.Add("secret/a", "env:prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, err := tr.Tags("secret/a")
	if err != nil {
		t.Fatalf("Tags error: %v", err)
	}
	if len(tags) != 1 || tags[0] != "env:prod" {
		t.Errorf("expected [env:prod], got %v", tags)
	}
}

func TestAdd_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTagger()
	if err := tr.Add("", "tag"); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestAdd_EmptyTag_ReturnsError(t *testing.T) {
	tr := newTagger()
	if err := tr.Add("secret/a", ""); err == nil {
		t.Error("expected error for empty tag")
	}
}

func TestAdd_Duplicate_Ignored(t *testing.T) {
	tr := newTagger()
	_ = tr.Add("secret/a", "env:prod")
	_ = tr.Add("secret/a", "env:prod")
	tags, _ := tr.Tags("secret/a")
	if len(tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(tags))
	}
}

func TestRemove_ExistingTag(t *testing.T) {
	tr := newTagger()
	_ = tr.Add("secret/a", "env:prod")
	if err := tr.Remove("secret/a", "env:prod"); err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	if _, err := tr.Tags("secret/a"); err == nil {
		t.Error("expected not-found error after removing last tag")
	}
}

func TestRemove_UnknownPath_ReturnsError(t *testing.T) {
	tr := newTagger()
	if err := tr.Remove("secret/missing", "tag"); err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestPathsWithTag_ReturnsMatchingPaths(t *testing.T) {
	tr := newTagger()
	_ = tr.Add("secret/a", "team:infra")
	_ = tr.Add("secret/b", "team:infra")
	_ = tr.Add("secret/c", "team:dev")
	paths := tr.PathsWithTag("team:infra")
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	tr := newTagger()
	_ = tr.Add("secret/a", "env:prod")
	out := secrettag.FormatTable(tr, []string{"secret/a"})
	if !strings.Contains(out, "PATH") || !strings.Contains(out, "TAGS") {
		t.Errorf("expected headers in output, got: %s", out)
	}
	if !strings.Contains(out, "env:prod") {
		t.Errorf("expected tag value in output, got: %s", out)
	}
}

func TestFormatTable_EmptyPaths(t *testing.T) {
	tr := newTagger()
	out := secrettag.FormatTable(tr, nil)
	if !strings.Contains(out, "No tagged") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	tr := newTagger()
	_ = tr.Add("secret/a", "x")
	_ = tr.Add("secret/b", "y")
	out := secrettag.FormatSummary(tr)
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in summary, got: %s", out)
	}
}
