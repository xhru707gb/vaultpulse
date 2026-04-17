package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/baseline"
)

var fixed = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makeEntry(path string, version int, ttl int64) baseline.Entry {
	return baseline.Entry{
		Path:        path,
		Version:     version,
		TTL:         ttl,
		LastRotated: fixed,
	}
}

func TestCompare_NoDrift(t *testing.T) {
	s := baseline.New()
	e := makeEntry("secret/db", 2, 3600)
	s.Record(e)
	drifts, err := s.Compare(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(drifts) != 0 {
		t.Fatalf("expected no drift, got %d", len(drifts))
	}
}

func TestCompare_VersionDrift(t *testing.T) {
	s := baseline.New()
	s.Record(makeEntry("secret/db", 2, 3600))
	current := makeEntry("secret/db", 3, 3600)
	drifts, err := s.Compare(current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(drifts) != 1 || drifts[0].Field != "version" {
		t.Fatalf("expected version drift, got %+v", drifts)
	}
}

func TestCompare_MultipleDrifts(t *testing.T) {
	s := baseline.New()
	s.Record(makeEntry("secret/api", 1, 7200))
	current := makeEntry("secret/api", 2, 3600)
	drifts, err := s.Compare(current)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(drifts) != 2 {
		t.Fatalf("expected 2 drifts, got %d", len(drifts))
	}
}

func TestCompare_UnknownPath(t *testing.T) {
	s := baseline.New()
	_, err := s.Compare(makeEntry("secret/missing", 1, 0))
	if err == nil {
		t.Fatal("expected error for unknown path")
	}
}

func TestSaveAndLoadJSON(t *testing.T) {
	s := baseline.New()
	s.Record(makeEntry("secret/db", 5, 1800))
	tmp := filepath.Join(t.TempDir(), "baseline.json")
	if err := s.SaveJSON(tmp); err != nil {
		t.Fatalf("SaveJSON: %v", err)
	}
	s2 := baseline.New()
	if err := s2.LoadJSON(tmp); err != nil {
		t.Fatalf("LoadJSON: %v", err)
	}
	drifts, err := s2.Compare(makeEntry("secret/db", 5, 1800))
	if err != nil {
		t.Fatalf("Compare after load: %v", err)
	}
	if len(drifts) != 0 {
		t.Fatalf("expected no drift after round-trip, got %d", len(drifts))
	}
}

func TestLoadJSON_FileNotFound(t *testing.T) {
	s := baseline.New()
	err := s.LoadJSON("/nonexistent/baseline.json")
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}

func TestFormatDrifts_NoDrift(t *testing.T) {
	out := baseline.FormatDrifts(nil)
	if out != "No drift detected.\n" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestFormatDrifts_ContainsHeaders(t *testing.T) {
	drifts := []baseline.Drift{{Path: "secret/db", Field: "version", Was: "1", Now: "2"}}
	out := baseline.FormatDrifts(drifts)
	for _, h := range []string{"PATH", "FIELD", "WAS", "NOW"} {
		if !contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
