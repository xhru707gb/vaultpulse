package ownership_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpulse/internal/ownership"
)

func newRegistry() *ownership.Registry {
	return ownership.New()
}

func TestRegister_And_Get(t *testing.T) {
	r := newRegistry()
	err := r.Register("secret/db", "alice", "platform", "alice@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := r.Get("secret/db")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Owner != "alice" {
		t.Errorf("expected owner alice, got %s", e.Owner)
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Register("", "alice", "team", "c"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRegister_EmptyOwner_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Register("secret/x", "", "team", "c"); err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestGet_UnknownPath_ReturnsFalse(t *testing.T) {
	r := newRegistry()
	_, ok := r.Get("secret/missing")
	if ok {
		t.Fatal("expected false for unknown path")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	r := newRegistry()
	_ = r.Register("secret/a", "bob", "ops", "bob@example.com")
	if err := r.Remove("secret/a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := r.Get("secret/a")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	r := newRegistry()
	if err := r.Remove("secret/nope"); err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	r := newRegistry()
	_ = r.Register("secret/db", "alice", "platform", "alice@example.com")
	out := ownership.FormatTable(r.All())
	for _, h := range []string{"PATH", "OWNER", "TEAM", "CONTACT"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := ownership.FormatTable(nil)
	if !strings.Contains(out, "no ownership") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
