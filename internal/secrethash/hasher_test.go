package secrethash_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secrethash"
)

func newTestHasher() *secrethash.Hasher {
	return secrethash.New()
}

func TestRecord_NewEntry_ReturnsChanged(t *testing.T) {
	h := newTestHasher()
	changed, err := h.Record("secret/foo", "mysecretvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for new entry")
	}
}

func TestRecord_SameValue_ReturnsNotChanged(t *testing.T) {
	h := newTestHasher()
	h.Record("secret/foo", "mysecretvalue")
	changed, err := h.Record("secret/foo", "mysecretvalue")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false for same value")
	}
}

func TestRecord_DifferentValue_ReturnsChanged(t *testing.T) {
	h := newTestHasher()
	h.Record("secret/foo", "value1")
	changed, err := h.Record("secret/foo", "value2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for different value")
	}
}

func TestRecord_VersionIncrements(t *testing.T) {
	h := newTestHasher()
	h.Record("secret/foo", "v1")
	h.Record("secret/foo", "v2")
	e, ok := h.Get("secret/foo")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Version != 2 {
		t.Errorf("expected version 2, got %d", e.Version)
	}
}

func TestRecord_EmptyPath_ReturnsError(t *testing.T) {
	h := newTestHasher()
	_, err := h.Record("", "value")
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRecord_EmptyValue_ReturnsError(t *testing.T) {
	h := newTestHasher()
	_, err := h.Record("secret/foo", "")
	if err == nil {
		t.Error("expected error for empty value")
	}
}

func TestGet_NotFound_ReturnsFalse(t *testing.T) {
	h := newTestHasher()
	_, ok := h.Get("secret/missing")
	if ok {
		t.Error("expected ok=false for unknown path")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	h := newTestHasher()
	h.Record("secret/a", "val1")
	h.Record("secret/b", "val2")
	all := h.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	h := newTestHasher()
	h.Record("secret/demo", "somevalue")
	out := secrethash.FormatTable(h.All())
	for _, hdr := range []string{"PATH", "HASH", "VERSION", "CHANGED AT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := secrethash.FormatTable(nil)
	if !strings.Contains(out, "No hash entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Count(t *testing.T) {
	entries := []secrethash.Entry{
		{Path: "a", Hash: "abc", ChangedAt: time.Now(), Version: 1},
		{Path: "b", Hash: "def", ChangedAt: time.Now(), Version: 1},
	}
	out := secrethash.FormatSummary(entries)
	if !strings.Contains(out, "2") {
		t.Errorf("expected count 2 in summary, got: %s", out)
	}
}
