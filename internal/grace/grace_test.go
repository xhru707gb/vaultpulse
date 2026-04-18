package grace

import (
	"strings"
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTracker(t *testing.T) *Tracker {
	t.Helper()
	tr, err := New(2*time.Hour, func() time.Time { return fixedNow })
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(-1*time.Second, nil)
	if err == nil {
		t.Fatal("expected error for non-positive window")
	}
}

func TestRegister_And_Active(t *testing.T) {
	tr := newTracker(t)
	expiredAt := fixedNow.Add(-30 * time.Minute) // expired 30m ago, grace ends in 90m
	if err := tr.Register("secret/a", expiredAt); err != nil {
		t.Fatalf("Register: %v", err)
	}
	actives := tr.Active()
	if len(actives) != 1 {
		t.Fatalf("expected 1 active, got %d", len(actives))
	}
	if actives[0].Path != "secret/a" {
		t.Errorf("unexpected path: %s", actives[0].Path)
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	expiredAt := fixedNow.Add(-10 * time.Minute)
	_ = tr.Register("secret/b", expiredAt)
	if err := tr.Register("secret/b", expiredAt); err != ErrAlreadyTracked {
		t.Errorf("expected ErrAlreadyTracked, got %v", err)
	}
}

func TestActive_ExcludesExpiredGrace(t *testing.T) {
	tr := newTracker(t)
	// expired 3 hours ago — grace window of 2h already passed
	_ = tr.Register("secret/old", fixedNow.Add(-3*time.Hour))
	if len(tr.Active()) != 0 {
		t.Error("expected no active entries past grace window")
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	tr := newTracker(t)
	_ = tr.Register("secret/c", fixedNow.Add(-10*time.Minute))
	if err := tr.Remove("secret/c"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if len(tr.Active()) != 0 {
		t.Error("expected empty after remove")
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	tr := newTracker(t)
	if err := tr.Remove("secret/x"); err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	entry := Entry{
		Path:        "secret/demo",
		ExpiredAt:   fixedNow.Add(-30 * time.Minute),
		GraceEndsAt: fixedNow.Add(90 * time.Minute),
	}
	out := FormatTable([]Entry{entry}, fixedNow)
	for _, h := range []string{colPath, colExpired, colGraceEnd, colRemains} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := FormatTable(nil, fixedNow)
	if !strings.Contains(out, "No secrets") {
		t.Error("expected empty message")
	}
}
