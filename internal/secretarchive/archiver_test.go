package secretarchive_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretarchive"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func newTestArchiver() *secretarchive.Archiver {
	a := secretarchive.New()
	return a
}

func TestArchive_And_Len(t *testing.T) {
	a := newTestArchiver()
	if err := a.Archive("secret/db", 3, "rotated"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := a.Len(); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestArchive_EmptyPath_ReturnsError(t *testing.T) {
	a := newTestArchiver()
	if err := a.Archive("", 1, "rotated"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestArchive_InvalidVersion_ReturnsError(t *testing.T) {
	a := newTestArchiver()
	if err := a.Archive("secret/db", 0, "rotated"); err == nil {
		t.Fatal("expected error for version < 1")
	}
}

func TestArchive_EmptyReason_ReturnsError(t *testing.T) {
	a := newTestArchiver()
	if err := a.Archive("secret/db", 1, ""); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestForPath_ReturnsMatchingEntries(t *testing.T) {
	a := newTestArchiver()
	_ = a.Archive("secret/db", 1, "expired")
	_ = a.Archive("secret/db", 2, "rotated")
	_ = a.Archive("secret/api", 1, "expired")

	entries := a.ForPath("secret/db")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for secret/db, got %d", len(entries))
	}
}

func TestForPath_UnknownPath_ReturnsEmpty(t *testing.T) {
	a := newTestArchiver()
	_ = a.Archive("secret/db", 1, "expired")

	entries := a.ForPath("secret/unknown")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	a := newTestArchiver()
	_ = a.Archive("secret/db", 1, "rotated")
	all := a.All()
	all[0].Path = "mutated"

	original := a.All()
	if original[0].Path == "mutated" {
		t.Fatal("All() should return a copy, not a reference")
	}
}

func TestArchive_TimestampIsUTC(t *testing.T) {
	a := newTestArchiver()
	_ = a.Archive("secret/db", 1, "rotated")
	entries := a.All()
	if entries[0].ArchivedAt.Location() != time.UTC {
		t.Fatal("expected UTC timestamp")
	}
}
