package tagindex_test

import (
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/tagindex"
)

func newIndex(t *testing.T) *tagindex.Index {
	t.Helper()
	return tagindex.New()
}

func TestAdd_And_Paths(t *testing.T) {
	idx := newIndex(t)
	if err := idx.Add("env:prod", "secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths := idx.Paths("env:prod")
	if len(paths) != 1 || paths[0] != "secret/db" {
		t.Fatalf("expected [secret/db], got %v", paths)
	}
}

func TestAdd_EmptyTag_ReturnsError(t *testing.T) {
	idx := newIndex(t)
	if err := idx.Add("", "secret/db"); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestAdd_EmptyPath_ReturnsError(t *testing.T) {
	idx := newIndex(t)
	if err := idx.Add("env:prod", ""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	idx := newIndex(t)
	_ = idx.Add("env:prod", "secret/db")
	if err := idx.Remove("env:prod", "secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Len("env:prod") != 0 {
		t.Fatal("expected tag to be empty after removal")
	}
}

func TestRemove_CleansUpEmptyTag(t *testing.T) {
	idx := newIndex(t)
	_ = idx.Add("team:sre", "secret/x")
	_ = idx.Remove("team:sre", "secret/x")
	tags := idx.Tags()
	for _, tag := range tags {
		if tag == "team:sre" {
			t.Fatal("empty tag should be removed from index")
		}
	}
}

func TestPaths_SortedOrder(t *testing.T) {
	idx := newIndex(t)
	_ = idx.Add("env:prod", "secret/z")
	_ = idx.Add("env:prod", "secret/a")
	_ = idx.Add("env:prod", "secret/m")
	paths := idx.Paths("env:prod")
	if paths[0] != "secret/a" || paths[1] != "secret/m" || paths[2] != "secret/z" {
		t.Fatalf("expected sorted paths, got %v", paths)
	}
}

func TestTags_SortedOrder(t *testing.T) {
	idx := newIndex(t)
	_ = idx.Add("zzz", "secret/a")
	_ = idx.Add("aaa", "secret/b")
	tags := idx.Tags()
	if tags[0] != "aaa" || tags[1] != "zzz" {
		t.Fatalf("expected sorted tags, got %v", tags)
	}
}

func TestLen_UnknownTag_ReturnsZero(t *testing.T) {
	idx := newIndex(t)
	if idx.Len("nonexistent") != 0 {
		t.Fatal("expected 0 for unknown tag")
	}
}
