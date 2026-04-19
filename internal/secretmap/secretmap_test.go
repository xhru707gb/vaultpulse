package secretmap_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretmap"
)

func makeEntry(path string) secretmap.Entry {
	return secretmap.Entry{
		Path:    path,
		Version: 1,
		Owner:   "team-a",
	}
}

func TestRegister_And_Get(t *testing.T) {
	r := secretmap.New()
	e := makeEntry("secret/foo")
	if err := r.Register(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := r.Get("secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != e.Path {
		t.Errorf("expected %s, got %s", e.Path, got.Path)
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	r := secretmap.New()
	_ = r.Register(makeEntry("secret/foo"))
	err := r.Register(makeEntry("secret/foo"))
	if err != secretmap.ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	r := secretmap.New()
	_, err := r.Get("secret/missing")
	if err != secretmap.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	r := secretmap.New()
	_ = r.Register(makeEntry("secret/bar"))
	if err := r.Remove("secret/bar"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", r.Len())
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	r := secretmap.New()
	err := r.Remove("secret/ghost")
	if err != secretmap.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	r := secretmap.New()
	_ = r.Register(secretmap.Entry{Path: "secret/x", Version: 2, Owner: "ops"})
	out := secretmap.FormatTable(r.All(), time.Now())
	for _, h := range []string{"PATH", "VERSION", "EXPIRES", "OWNER"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_EmptyRegistry(t *testing.T) {
	out := secretmap.FormatTable(nil, time.Now())
	if !strings.Contains(out, "no secrets") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
