package diff_test

import (
	"testing"
	"time"

	"github.com/vaultpulse/vaultpulse/internal/diff"
)

var baseTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func entry(path string, version int, expiresAt time.Time) diff.SecretEntry {
	return diff.SecretEntry{Path: path, Version: version, ExpiresAt: expiresAt}
}

func TestCompute_AddedEntry(t *testing.T) {
	prev := []diff.SecretEntry{}
	curr := []diff.SecretEntry{entry("secret/new", 1, baseTime)}

	changes := diff.Compute(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.ChangeAdded {
		t.Errorf("expected added, got %s", changes[0].Kind)
	}
	if changes[0].Path != "secret/new" {
		t.Errorf("unexpected path: %s", changes[0].Path)
	}
}

func TestCompute_RemovedEntry(t *testing.T) {
	prev := []diff.SecretEntry{entry("secret/gone", 2, baseTime)}
	curr := []diff.SecretEntry{}

	changes := diff.Compute(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.ChangeRemoved {
		t.Errorf("expected removed, got %s", changes[0].Kind)
	}
}

func TestCompute_ModifiedVersion(t *testing.T) {
	prev := []diff.SecretEntry{entry("secret/a", 1, baseTime)}
	curr := []diff.SecretEntry{entry("secret/a", 2, baseTime)}

	changes := diff.Compute(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != diff.ChangeModified {
		t.Errorf("expected modified, got %s", changes[0].Kind)
	}
}

func TestCompute_ModifiedExpiry(t *testing.T) {
	prev := []diff.SecretEntry{entry("secret/b", 1, baseTime)}
	curr := []diff.SecretEntry{entry("secret/b", 1, baseTime.Add(24*time.Hour))}

	changes := diff.Compute(prev, curr)
	if len(changes) != 1 || changes[0].Kind != diff.ChangeModified {
		t.Errorf("expected modified due to expiry change")
	}
}

func TestCompute_NoChanges(t *testing.T) {
	prev := []diff.SecretEntry{entry("secret/stable", 3, baseTime)}
	curr := []diff.SecretEntry{entry("secret/stable", 3, baseTime)}

	changes := diff.Compute(prev, curr)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestChange_String_Added(t *testing.T) {
	c := diff.Change{
		Path: "secret/x",
		Kind: diff.ChangeAdded,
		Curr: &diff.SecretEntry{Path: "secret/x", Version: 1},
	}
	got := c.String()
	if got == "" {
		t.Error("expected non-empty string")
	}
}
