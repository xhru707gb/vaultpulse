package acknowledge

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func newTestTracker(window time.Duration) *Tracker {
	t, _ := New(window)
	t.now = func() time.Time { return fixedNow }
	return t
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestAcknowledge_And_IsAcknowledged(t *testing.T) {
	tr := newTestTracker(time.Hour)
	if err := tr.Acknowledge("secret/foo", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tr.IsAcknowledged("secret/foo") {
		t.Fatal("expected path to be acknowledged")
	}
}

func TestAcknowledge_Duplicate_ReturnsError(t *testing.T) {
	tr := newTestTracker(time.Hour)
	_ = tr.Acknowledge("secret/foo", "alice")
	err := tr.Acknowledge("secret/foo", "bob")
	if err != ErrAlreadyAcknowledged {
		t.Fatalf("expected ErrAlreadyAcknowledged, got %v", err)
	}
}

func TestIsAcknowledged_AfterExpiry_ReturnsFalse(t *testing.T) {
	tr := newTestTracker(time.Minute)
	_ = tr.Acknowledge("secret/bar", "alice")
	tr.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }
	if tr.IsAcknowledged("secret/bar") {
		t.Fatal("expected acknowledgement to be expired")
	}
}

func TestRevoke_RemovesEntry(t *testing.T) {
	tr := newTestTracker(time.Hour)
	_ = tr.Acknowledge("secret/baz", "alice")
	if err := tr.Revoke("secret/baz"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.IsAcknowledged("secret/baz") {
		t.Fatal("expected path to no longer be acknowledged")
	}
}

func TestRevoke_Unknown_ReturnsError(t *testing.T) {
	tr := newTestTracker(time.Hour)
	if err := tr.Revoke("secret/unknown"); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestList_ReturnsActiveOnly(t *testing.T) {
	tr := newTestTracker(time.Minute)
	_ = tr.Acknowledge("secret/a", "alice")
	_ = tr.Acknowledge("secret/b", "bob")
	tr.now = func() time.Time { return fixedNow.Add(2 * time.Minute) }
	// re-ack secret/a under new time
	tr.entries["secret/a"] = Entry{
		Path: "secret/a", AckedAt: fixedNow.Add(2 * time.Minute),
		ExpiresAt: fixedNow.Add(3 * time.Minute), AckedBy: "alice",
	}
	list := tr.List()
	if len(list) != 1 || list[0].Path != "secret/a" {
		t.Fatalf("expected 1 active entry, got %d", len(list))
	}
}
