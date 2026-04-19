package secretgroup_test

import (
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretgroup"
)

func newGrouper() *secretgroup.Grouper {
	return secretgroup.New()
}

func TestAdd_And_Get(t *testing.T) {
	g := newGrouper()
	if err := g.Add("infra", "secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	grp, ok := g.Get("infra")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if len(grp.Paths) != 1 || grp.Paths[0] != "secret/db" {
		t.Errorf("unexpected paths: %v", grp.Paths)
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	g := newGrouper()
	if err := g.Add("", "secret/db"); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_EmptyPath_ReturnsError(t *testing.T) {
	g := newGrouper()
	if err := g.Add("infra", ""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAdd_Duplicate_NoError_NoDuplicate(t *testing.T) {
	g := newGrouper()
	_ = g.Add("infra", "secret/db")
	_ = g.Add("infra", "secret/db")
	grp, _ := g.Get("infra")
	if len(grp.Paths) != 1 {
		t.Errorf("expected 1 path, got %d", len(grp.Paths))
	}
}

func TestRemove_ExistingPath(t *testing.T) {
	g := newGrouper()
	_ = g.Add("infra", "secret/db")
	if err := g.Remove("infra", "secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	grp, _ := g.Get("infra")
	if len(grp.Paths) != 0 {
		t.Errorf("expected empty paths, got %v", grp.Paths)
	}
}

func TestRemove_UnknownGroup_ReturnsError(t *testing.T) {
	g := newGrouper()
	if err := g.Remove("missing", "secret/db"); err == nil {
		t.Fatal("expected error")
	}
}

func TestAll_SortedByName(t *testing.T) {
	g := newGrouper()
	_ = g.Add("z-group", "secret/z")
	_ = g.Add("a-group", "secret/a")
	all := g.All()
	if len(all) != 2 || all[0].Name != "a-group" {
		t.Errorf("unexpected order: %v", all)
	}
}

func TestFindByPrefix_MatchesGroup(t *testing.T) {
	g := newGrouper()
	_ = g.Add("infra", "secret/db/prod")
	_ = g.Add("app", "secret/app/config")
	result := g.FindByPrefix("secret/db")
	if len(result) != 1 || result[0].Name != "infra" {
		t.Errorf("unexpected result: %v", result)
	}
}
