package quarantine_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/quarantine"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestAdd_NewEntry(t *testing.T) {
	s := quarantine.New()
	err := s.Add("secret/db", quarantine.ReasonLeaked, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsQuarantined("secret/db") {
		t.Fatal("expected path to be quarantined")
	}
}

func TestAdd_Duplicate_ReturnsError(t *testing.T) {
	s := quarantine.New()
	_ = s.Add("secret/db", quarantine.ReasonExpired, "")
	err := s.Add("secret/db", quarantine.ReasonManual, "")
	if err != quarantine.ErrAlreadyQuarantined {
		t.Fatalf("expected ErrAlreadyQuarantined, got %v", err)
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	s := quarantine.New()
	_ = s.Add("secret/x", quarantine.ReasonPolicy, "")
	if err := s.Remove("secret/x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsQuarantined("secret/x") {
		t.Fatal("expected path to be removed")
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	s := quarantine.New()
	err := s.Remove("secret/missing")
	if err != quarantine.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := quarantine.New()
	_ = s.Add("secret/a", quarantine.ReasonLeaked, "note a")
	_ = s.Add("secret/b", quarantine.ReasonExpired, "note b")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	s := quarantine.New()
	_ = s.Add("secret/api", quarantine.ReasonManual, "manual flag")
	out := quarantine.FormatTable(s.All())
	for _, hdr := range []string{"PATH", "REASON", "QUARANTINED AT", "NOTE"} {
		if !containsStr(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_EmptyEntries(t *testing.T) {
	out := quarantine.FormatTable(nil)
	if !containsStr(out, "no quarantined") {
		t.Error("expected empty message")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
