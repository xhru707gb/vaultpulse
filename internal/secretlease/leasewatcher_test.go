package secretlease

import (
	"testing"
	"time"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestWatcher(t *testing.T) *Watcher {
	t.Helper()
	w, err := New(10 * time.Minute)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	w.now = func() time.Time { return fixedNow }
	return w
}

func TestNew_InvalidWarnBefore(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero warnBefore")
	}
}

func TestRegister_And_Evaluate_OK(t *testing.T) {
	w := newTestWatcher(t)
	err := w.Register(Entry{
		Path:      "secret/db",
		LeaseID:   "lease-1",
		ExpiresAt: fixedNow.Add(30 * time.Minute),
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	statuses := w.Evaluate()
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].State != StateOK {
		t.Errorf("expected StateOK, got %v", statuses[0].State)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	w := newTestWatcher(t)
	_ = w.Register(Entry{
		Path:      "secret/api",
		LeaseID:   "lease-2",
		ExpiresAt: fixedNow.Add(5 * time.Minute), // within 10m warn threshold
	})
	statuses := w.Evaluate()
	if statuses[0].State != StateWarning {
		t.Errorf("expected StateWarning, got %v", statuses[0].State)
	}
}

func TestEvaluate_Expired(t *testing.T) {
	w := newTestWatcher(t)
	_ = w.Register(Entry{
		Path:      "secret/old",
		LeaseID:   "lease-3",
		ExpiresAt: fixedNow.Add(-1 * time.Minute),
	})
	statuses := w.Evaluate()
	if statuses[0].State != StateExpired {
		t.Errorf("expected StateExpired, got %v", statuses[0].State)
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	w := newTestWatcher(t)
	err := w.Register(Entry{LeaseID: "x", ExpiresAt: fixedNow.Add(time.Hour)})
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRegister_EmptyLeaseID_ReturnsError(t *testing.T) {
	w := newTestWatcher(t)
	err := w.Register(Entry{Path: "secret/x", ExpiresAt: fixedNow.Add(time.Hour)})
	if err == nil {
		t.Fatal("expected error for empty leaseID")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	w := newTestWatcher(t)
	_ = w.Register(Entry{Path: "secret/tmp", LeaseID: "l", ExpiresAt: fixedNow.Add(time.Hour)})
	w.Remove("secret/tmp")
	if len(w.Evaluate()) != 0 {
		t.Error("expected no entries after Remove")
	}
}
