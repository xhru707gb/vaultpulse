package trend_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/trend"
)

var base = time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)

func makeEntries() []trend.Entry {
	return []trend.Entry{
		{Path: "secret/db", EventType: "expiry", Timestamp: base},
		{Path: "secret/db", EventType: "expiry", Timestamp: base.Add(30 * time.Minute)},
		{Path: "secret/db", EventType: "expiry", Timestamp: base.Add(25 * time.Hour)},
		{Path: "secret/api", EventType: "rotation", Timestamp: base},
	}
}

func TestNewAnalyzer_InvalidBucket(t *testing.T) {
	_, err := trend.NewAnalyzer(0)
	if err == nil {
		t.Fatal("expected error for zero bucket size")
	}
}

func TestAnalyse_GroupsByPathAndEvent(t *testing.T) {
	a, _ := trend.NewAnalyzer(24 * time.Hour)
	reports := a.Analyse(makeEntries())
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports, got %d", len(reports))
	}
}

func TestAnalyse_TotalCount(t *testing.T) {
	a, _ := trend.NewAnalyzer(24 * time.Hour)
	reports := a.Analyse(makeEntries())
	var dbReport *trend.Report
	for i := range reports {
		if reports[i].Path == "secret/db" {
			dbReport = &reports[i]
		}
	}
	if dbReport == nil {
		t.Fatal("secret/db report not found")
	}
	if dbReport.TotalCount != 3 {
		t.Errorf("expected total 3, got %d", dbReport.TotalCount)
	}
}

func TestAnalyse_BucketsCorrectly(t *testing.T) {
	a, _ := trend.NewAnalyzer(24 * time.Hour)
	reports := a.Analyse(makeEntries())
	for _, r := range reports {
		if r.Path == "secret/db" {
			if len(r.Points) != 2 {
				t.Errorf("expected 2 buckets for secret/db, got %d", len(r.Points))
			}
		}
	}
}

func TestAnalyse_EmptyInput(t *testing.T) {
	a, _ := trend.NewAnalyzer(time.Hour)
	reports := a.Analyse(nil)
	if len(reports) != 0 {
		t.Errorf("expected empty reports, got %d", len(reports))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	a, _ := trend.NewAnalyzer(24 * time.Hour)
	reports := a.Analyse(makeEntries())
	out := trend.FormatTable(reports)
	for _, h := range []string{"PATH", "EVENT", "TOTAL", "TREND"} {
		if !contains(out, h) {
			t.Errorf("output missing header %q", h)
		}
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := trend.FormatTable(nil)
	if !contains(out, "No trend") {
		t.Error("expected no-data message")
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
