package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/expiry"
	"github.com/yourusername/vaultpulse/internal/rotation"
	"github.com/yourusername/vaultpulse/internal/snapshot"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestWriter(t *testing.T) (*snapshot.Writer, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "snapshot.json")
	w := snapshot.NewWriter(path)
	w.SetNow(func() time.Time { return fixedNow })
	return w, path
}

func TestWrite_CreatesFile(t *testing.T) {
	w, path := newTestWriter(t)

	err := w.Write([]expiry.Status{{Path: "secret/a"}}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestWrite_RoundTrip(t *testing.T) {
	w, path := newTestWriter(t)

	es := []expiry.Status{{Path: "secret/db", TTL: 3600}}
	rs := []rotation.Status{{Path: "secret/api", Overdue: true}}

	if err := w.Write(es, rs); err != nil {
		t.Fatalf("write: %v", err)
	}

	snap, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if snap.CapturedAt != fixedNow {
		t.Errorf("captured_at: got %v, want %v", snap.CapturedAt, fixedNow)
	}
	if len(snap.ExpiryStatuses) != 1 || snap.ExpiryStatuses[0].Path != "secret/db" {
		t.Errorf("expiry statuses mismatch")
	}
	if len(snap.RotationStatuses) != 1 || !snap.RotationStatuses[0].Overdue {
		t.Errorf("rotation statuses mismatch")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
