package secretpin_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretpin"
)

func newTestPinner() *secretpin.Pinner {
	return secretpin.New()
}

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestPin_And_Check_NoDrift(t *testing.T) {
	p := newTestPinner()
	if err := p.Pin("secret/db", 3, "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, err := p.Check("secret/db", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Drifted {
		t.Error("expected no drift")
	}
}

func TestPin_And_Check_Drift(t *testing.T) {
	p := newTestPinner()
	_ = p.Pin("secret/db", 3, "alice")
	res, err := p.Check("secret/db", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Drifted {
		t.Error("expected drift")
	}
	if res.PinnedVersion != 3 || res.CurrentVersion != 5 {
		t.Errorf("unexpected versions: pinned=%d current=%d", res.PinnedVersion, res.CurrentVersion)
	}
}

func TestPin_Duplicate_ReturnsError(t *testing.T) {
	p := newTestPinner()
	_ = p.Pin("secret/db", 1, "alice")
	err := p.Pin("secret/db", 2, "bob")
	if err == nil {
		t.Fatal("expected error for duplicate pin")
	}
}

func TestUnpin_RemovesPin(t *testing.T) {
	p := newTestPinner()
	_ = p.Pin("secret/db", 1, "alice")
	if err := p.Unpin("secret/db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := p.Check("secret/db", 1)
	if err == nil {
		t.Error("expected error after unpin")
	}
}

func TestUnpin_Unknown_ReturnsError(t *testing.T) {
	p := newTestPinner()
	err := p.Unpin("secret/unknown")
	if err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestPin_InvalidArgs_ReturnsError(t *testing.T) {
	p := newTestPinner()
	if err := p.Pin("", 1, "alice"); err == nil {
		t.Error("expected error for empty path")
	}
	if err := p.Pin("secret/db", 0, "alice"); err == nil {
		t.Error("expected error for zero version")
	}
	if err := p.Pin("secret/db", 1, ""); err == nil {
		t.Error("expected error for empty pinnedBy")
	}
}

func TestFormatDrifts_ContainsHeaders(t *testing.T) {
	results := []secretpin.DriftResult{
		{Path: "secret/db", PinnedVersion: 2, CurrentVersion: 4, Drifted: true},
		{Path: "secret/api", PinnedVersion: 1, CurrentVersion: 1, Drifted: false},
	}
	out := secretpin.FormatDrifts(results)
	for _, want := range []string{"PATH", "PINNED", "CURRENT", "DRIFTED", "secret/db", "YES", "no"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q", want)
		}
	}
}

func TestFormatPins_ContainsEntry(t *testing.T) {
	p := newTestPinner()
	_ = p.Pin("secret/db", 3, "alice")
	out := secretpin.FormatPins(p.All())
	if !strings.Contains(out, "secret/db") {
		t.Error("expected output to contain path")
	}
	if !strings.Contains(out, "alice") {
		t.Error("expected output to contain pinnedBy")
	}
}
