package lease_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/lease"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTracker() *lease.Tracker {
	return lease.New(func() time.Time { return fixedNow })
}

func makeEntry(id, path string, ttl time.Duration) lease.Entry {
	return lease.Entry{
		LeaseID:   id,
		Path:      path,
		IssuedAt:  fixedNow,
		ExpiresAt: fixedNow.Add(ttl),
		Renewable: true,
	}
}

func TestRegister_And_Get(t *testing.T) {
	tr := newTracker()
	e := makeEntry("lease-1", "secret/db", 10*time.Minute)
	tr.Register(e)

	got, err := tr.Get("lease-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Path != "secret/db" {
		t.Errorf("expected path secret/db, got %s", got.Path)
	}
}

func TestGet_NotFound(t *testing.T) {
	tr := newTracker()
	_, err := tr.Get("missing")
	if err != lease.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestRemove_ExistingEntry(t *testing.T) {
	tr := newTracker()
	tr.Register(makeEntry("lease-2", "secret/app", 5*time.Minute))
	if err := tr.Remove("lease-2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", tr.Len())
	}
}

func TestRemove_Unknown_ReturnsError(t *testing.T) {
	tr := newTracker()
	if err := tr.Remove("ghost"); err != lease.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestExpiring_FiltersCorrectly(t *testing.T) {
	tr := newTracker()
	tr.Register(makeEntry("soon", "secret/a", 2*time.Minute))
	tr.Register(makeEntry("later", "secret/b", 30*time.Minute))

	expiring := tr.Expiring(5 * time.Minute)
	if len(expiring) != 1 {
		t.Fatalf("expected 1 expiring lease, got %d", len(expiring))
	}
	if expiring[0].LeaseID != "soon" {
		t.Errorf("expected lease 'soon', got %s", expiring[0].LeaseID)
	}
}

func TestTTL_IsPositive(t *testing.T) {
	e := makeEntry("x", "secret/x", 10*time.Minute)
	ttl := e.TTL(fixedNow)
	if ttl != 10*time.Minute {
		t.Errorf("expected 10m TTL, got %v", ttl)
	}
}
