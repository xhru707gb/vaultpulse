package secretreport_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretreport"
)

func makeEntries() []secretreport.Entry {
	return []secretreport.Entry{
		{Path: "secret/a", Severity: secretreport.SeverityOK, TTL: 24 * time.Hour},
		{Path: "secret/b", Severity: secretreport.SeverityWarning, TTL: 2 * time.Hour},
		{Path: "secret/c", Severity: secretreport.SeverityCritical, TTL: 10 * time.Minute},
	}
}

func TestBuild_TotalCount(t *testing.T) {
	r := secretreport.New()
	rep, err := r.Build(makeEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.Total != 3 {
		t.Errorf("expected total 3, got %d", rep.Total)
	}
}

func TestBuild_AlertCount(t *testing.T) {
	r := secretreport.New()
	rep, err := r.Build(makeEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// warning + critical = 2
	if rep.AlertCount != 2 {
		t.Errorf("expected alert count 2, got %d", rep.AlertCount)
	}
}

func TestBuild_AllOK_ZeroAlerts(t *testing.T) {
	entries := []secretreport.Entry{
		{Path: "secret/x", Severity: secretreport.SeverityOK},
		{Path: "secret/y", Severity: secretreport.SeverityOK},
	}
	r := secretreport.New()
	rep, err := r.Build(entries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.AlertCount != 0 {
		t.Errorf("expected 0 alerts, got %d", rep.AlertCount)
	}
}

func TestBuild_NilEntries_ReturnsError(t *testing.T) {
	r := secretreport.New()
	_, err := r.Build(nil)
	if err == nil {
		t.Fatal("expected error for nil entries")
	}
}

func TestBuild_GeneratedAtIsUTC(t *testing.T) {
	r := secretreport.New()
	rep, err := r.Build(makeEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.GeneratedAt.Location() != time.UTC {
		t.Errorf("expected UTC, got %v", rep.GeneratedAt.Location())
	}
}
