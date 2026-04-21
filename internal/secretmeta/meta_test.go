package secretmeta_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpulse/internal/secretmeta"
)

func newRegistry() *secretmeta.Registry {
	return secretmeta.New()
}

func TestSet_And_Get(t *testing.T) {
	r := newRegistry()
	if err := r.Set("secret/db", "owner", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, err := r.Get("secret/db")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if m["owner"] != "alice" {
		t.Errorf("expected alice, got %s", m["owner"])
	}
}

func TestSet_EmptyPath_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Set("", "key", "val"); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestSet_EmptyKey_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Set("secret/db", "", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestGet_NotFound(t *testing.T) {
	r := newRegistry()
	_, err := r.Get("secret/missing")
	if err == nil {
		t.Error("expected ErrNotFound")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	r := newRegistry()
	_ = r.Set("secret/db", "env", "prod")
	if err := r.Delete("secret/db"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := r.Get("secret/db"); err == nil {
		t.Error("expected path to be deleted")
	}
}

func TestDelete_Unknown_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Delete("secret/ghost"); err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestPaths_ReturnsAllPaths(t *testing.T) {
	r := newRegistry()
	_ = r.Set("secret/a", "k", "v")
	_ = r.Set("secret/b", "k", "v")
	paths := r.Paths()
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(paths))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	r := newRegistry()
	_ = r.Set("secret/db", "owner", "alice")
	out := secretmeta.FormatTable(r)
	for _, hdr := range []string{"PATH", "KEY", "VALUE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestFormatTable_ContainsEntries(t *testing.T) {
	r := newRegistry()
	_ = r.Set("secret/db", "owner", "alice")
	out := secretmeta.FormatTable(r)
	if !strings.Contains(out, "secret/db") {
		t.Error("expected path in table output")
	}
	if !strings.Contains(out, "alice") {
		t.Error("expected value in table output")
	}
}

func TestFormatTable_EmptyRegistry(t *testing.T) {
	r := newRegistry()
	out := secretmeta.FormatTable(r)
	if !strings.Contains(out, "No metadata") {
		t.Error("expected empty message")
	}
}
