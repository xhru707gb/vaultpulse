package secretversion_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/secretversion"
)

func newTestTracker() *secretversion.Tracker {
	return secretversion.New()
}

func TestRegister_And_Get(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Register("secret/db", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Get("secret/db")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Version != 3 {
		t.Errorf("expected version 3, got %d", e.Version)
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Register("", 1); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRegister_InvalidVersion_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Register("secret/db", 0); err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestRegister_UpdatesExisting(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/db", 1)
	created := time.Now()
	_ = tr.Register("secret/db", 2)
	e, _ := tr.Get("secret/db")
	if e.Version != 2 {
		t.Errorf("expected version 2, got %d", e.Version)
	}
	if e.UpdatedAt.Before(created) {
		t.Error("UpdatedAt should be >= created time")
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := newTestTracker()
	_, ok := tr.Get("secret/missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/db", 1)
	if err := tr.Remove("secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := tr.Get("secret/db")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	tr := newTestTracker()
	if err := tr.Remove("secret/missing"); err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := newTestTracker()
	_ = tr.Register("secret/a", 1)
	_ = tr.Register("secret/b", 2)
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
