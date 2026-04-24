package secretreview_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretreview"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestReviewer(t *testing.T) *secretreview.Reviewer {
	t.Helper()
	r, err := secretreview.New(func() time.Time { return fixedNow })
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestRegister_And_Evaluate_OK(t *testing.T) {
	r := newTestReviewer(t)
	last := fixedNow.Add(-24 * time.Hour)
	if err := r.Register("secret/foo", "alice", 48*time.Hour, last); err != nil {
		t.Fatalf("Register: %v", err)
	}
	entries := r.Evaluate()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != secretreview.StatusApproved {
		t.Errorf("expected approved, got %v", entries[0].Status)
	}
}

func TestRegister_Overdue(t *testing.T) {
	r := newTestReviewer(t)
	last := fixedNow.Add(-72 * time.Hour)
	if err := r.Register("secret/bar", "bob", 48*time.Hour, last); err != nil {
		t.Fatalf("Register: %v", err)
	}
	entries := r.Evaluate()
	if entries[0].Status != secretreview.StatusOverdue {
		t.Errorf("expected overdue, got %v", entries[0].Status)
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	r := newTestReviewer(t)
	if err := r.Register("", "alice", time.Hour, fixedNow); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRegister_EmptyReviewer_ReturnsError(t *testing.T) {
	r := newTestReviewer(t)
	if err := r.Register("secret/x", "", time.Hour, fixedNow); err == nil {
		t.Error("expected error for empty reviewer")
	}
}

func TestRegister_InvalidInterval_ReturnsError(t *testing.T) {
	r := newTestReviewer(t)
	if err := r.Register("secret/x", "alice", 0, fixedNow); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestApprove_ResetsNextReview(t *testing.T) {
	r := newTestReviewer(t)
	last := fixedNow.Add(-72 * time.Hour)
	_ = r.Register("secret/baz", "carol", 48*time.Hour, last)
	if err := r.Approve("secret/baz"); err != nil {
		t.Fatalf("Approve: %v", err)
	}
	entries := r.Evaluate()
	if entries[0].Status != secretreview.StatusApproved {
		t.Errorf("expected approved after approval, got %v", entries[0].Status)
	}
}

func TestApprove_UnknownPath_ReturnsError(t *testing.T) {
	r := newTestReviewer(t)
	if err := r.Approve("secret/unknown"); err == nil {
		t.Error("expected error for unknown path")
	}
}
