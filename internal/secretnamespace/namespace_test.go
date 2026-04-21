package secretnamespace_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretnamespace"
)

func newRegistry() *secretnamespace.Registry {
	return secretnamespace.New()
}

func TestAdd_And_Paths(t *testing.T) {
	r := newRegistry()
	if err := r.Add("team-a", "secret/db/password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths, err := r.Paths("team-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || paths[0] != "secret/db/password" {
		t.Errorf("expected path not found, got %v", paths)
	}
}

func TestAdd_EmptyNamespace_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Add("", "secret/db/password"); err == nil {
		t.Fatal("expected error for empty namespace")
	}
}

func TestAdd_EmptyPath_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Add("team-a", ""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestAdd_Duplicate_ReturnsError(t *testing.T) {
	r := newRegistry()
	_ = r.Add("team-a", "secret/db/password")
	if err := r.Add("team-a", "secret/db/password"); err == nil {
		t.Fatal("expected duplicate error")
	}
}

func TestPaths_UnknownNamespace_ReturnsError(t *testing.T) {
	r := newRegistry()
	_, err := r.Paths("unknown")
	if err == nil {
		t.Fatal("expected error for unknown namespace")
	}
}

func TestNamespaces_SortedOrder(t *testing.T) {
	r := newRegistry()
	_ = r.Add("zebra", "secret/z")
	_ = r.Add("alpha", "secret/a")
	_ = r.Add("mango", "secret/m")
	ns := r.Namespaces()
	if ns[0] != "alpha" || ns[1] != "mango" || ns[2] != "zebra" {
		t.Errorf("unexpected order: %v", ns)
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	r := newRegistry()
	_ = r.Add("team-a", "secret/db/password")
	if err := r.Remove("team-a", "secret/db/password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	paths, _ := r.Paths("team-a")
	if len(paths) != 0 {
		t.Errorf("expected empty paths after remove, got %v", paths)
	}
}

func TestRemove_UnknownNamespace_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Remove("unknown", "secret/x"); err == nil {
		t.Fatal("expected error for unknown namespace")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	r := newRegistry()
	_ = r.Add("team-a", "secret/db/password")
	out := secretnamespace.FormatTable(r)
	for _, h := range []string{"NAMESPACE", "PATHS", "SAMPLE PATH"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q in output", h)
		}
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	r := newRegistry()
	_ = r.Add("team-a", "secret/db/password")
	_ = r.Add("team-a", "secret/api/key")
	_ = r.Add("team-b", "secret/tls/cert")
	summary := secretnamespace.FormatSummary(r)
	if !strings.Contains(summary, "namespaces=2") {
		t.Errorf("expected namespaces=2 in summary, got %q", summary)
	}
	if !strings.Contains(summary, "total_paths=3") {
		t.Errorf("expected total_paths=3 in summary, got %q", summary)
	}
}
