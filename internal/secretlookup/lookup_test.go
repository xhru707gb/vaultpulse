package secretlookup_test

import (
	"strings"
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/secretlookup"
)

func newIndex() *secretlookup.Index {
	return secretlookup.New()
}

func TestAdd_And_Lookup(t *testing.T) {
	idx := newIndex()
	_ = idx.Add("secret/a", "fp1")
	_ = idx.Add("secret/b", "fp1")
	paths := idx.Lookup("fp1")
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
}

func TestLookup_UnknownFingerprint_ReturnsNil(t *testing.T) {
	idx := newIndex()
	if paths := idx.Lookup("nope"); paths != nil {
		t.Fatalf("expected nil, got %v", paths)
	}
}

func TestAdd_EmptyPath_ReturnsError(t *testing.T) {
	idx := newIndex()
	if err := idx.Add("", "fp1"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAdd_EmptyFingerprint_ReturnsError(t *testing.T) {
	idx := newIndex()
	if err := idx.Add("secret/a", ""); err == nil {
		t.Fatal("expected error for empty fingerprint")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	idx := newIndex()
	_ = idx.Add("secret/a", "fp1")
	_ = idx.Add("secret/b", "fp1")
	_ = idx.Remove("secret/a", "fp1")
	paths := idx.Lookup("fp1")
	if len(paths) != 1 || paths[0] != "secret/b" {
		t.Fatalf("unexpected paths after remove: %v", paths)
	}
}

func TestRemove_LastPath_DeletesFingerprint(t *testing.T) {
	idx := newIndex()
	_ = idx.Add("secret/a", "fp1")
	_ = idx.Remove("secret/a", "fp1")
	if paths := idx.Lookup("fp1"); paths != nil {
		t.Fatalf("expected fingerprint to be gone, got %v", paths)
	}
}

func TestDuplicates_ReturnsOnlyShared(t *testing.T) {
	idx := newIndex()
	_ = idx.Add("secret/a", "fp-shared")
	_ = idx.Add("secret/b", "fp-shared")
	_ = idx.Add("secret/c", "fp-unique")
	dupes := idx.Duplicates()
	if _, ok := dupes["fp-shared"]; !ok {
		t.Fatal("expected fp-shared in duplicates")
	}
	if _, ok := dupes["fp-unique"]; ok {
		t.Fatal("fp-unique should not appear in duplicates")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretlookup.FormatTable(map[string][]string{
		"abc123": {"secret/a", "secret/b"},
	})
	for _, h := range []string{"FINGERPRINT", "PATH", "SHARED"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_EmptyDuplicates(t *testing.T) {
	out := secretlookup.FormatTable(map[string][]string{})
	if !strings.Contains(out, "No duplicate") {
		t.Errorf("expected no-duplicate message, got: %s", out)
	}
}
