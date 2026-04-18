package renew_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/renew"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestRenewer(t *testing.T) *renew.Renewer {
	t.Helper()
	r, err := renew.New(0.75)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return r
}

func TestNew_InvalidThreshold(t *testing.T) {
	_, err := renew.New(0)
	if err == nil {
		t.Fatal("expected error for threshold=0")
	}
	_, err = renew.New(1)
	if err == nil {
		t.Fatal("expected error for threshold=1")
	}
}

func TestRegister_And_Due(t *testing.T) {
	r := newTestRenewer(t)
	ttl := 4 * time.Hour
	if err := r.Register("secret/foo", ttl); err != nil {
		t.Fatalf("Register: %v", err)
	}
	// Not due immediately (renewAt is 75% through TTL)
	due := r.Due()
	if len(due) != 0 {
		t.Fatalf("expected 0 due, got %d", len(due))
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	r := newTestRenewer(t)
	if err := r.Register("secret/foo", time.Hour); err != nil {
		t.Fatalf("first Register: %v", err)
	}
	if err := r.Register("secret/foo", time.Hour); err == nil {
		t.Fatal("expected error on duplicate Register")
	}
}

func TestRecordRenewal_UpdatesEntry(t *testing.T) {
	r := newTestRenewer(t)
	if err := r.Register("secret/bar", time.Hour); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := r.RecordRenewal("secret/bar", 2*time.Hour); err != nil {
		t.Fatalf("RecordRenewal: %v", err)
	}
}

func TestRecordRenewal_UnknownPath_ReturnsError(t *testing.T) {
	r := newTestRenewer(t)
	if err := r.RecordRenewal("secret/unknown", time.Hour); err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestRemove_StopsTracking(t *testing.T) {
	r := newTestRenewer(t)
	if err := r.Register("secret/baz", time.Hour); err != nil {
		t.Fatalf("Register: %v", err)
	}
	r.Remove("secret/baz")
	if err := r.Register("secret/baz", time.Hour); err != nil {
		t.Fatalf("re-Register after Remove: %v", err)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := renew.FormatTable(nil)
	if out == "" {
		t.Fatal("expected non-empty output")
	}
}
