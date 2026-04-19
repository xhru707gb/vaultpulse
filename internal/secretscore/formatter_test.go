package secretscore_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretscore"
)

func makeResults() []secretscore.Result {
	return []secretscore.Result{
		{Path: "secret/db", Score: 80, Level: secretscore.RiskCritical, Reason: "TTL critical"},
		{Path: "secret/api", Score: 50, Level: secretscore.RiskHigh, Reason: "TTL high"},
		{Path: "secret/cache", Score: 0, Level: secretscore.RiskLow, Reason: "TTL ok"},
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretscore.FormatTable(makeResults())
	for _, h := range []string{"PATH", "SCORE", "RISK", "REASON"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestFormatTable_ContainsRiskLevels(t *testing.T) {
	out := secretscore.FormatTable(makeResults())
	for _, lvl := range []string{secretscore.RiskCritical, secretscore.RiskHigh, secretscore.RiskLow} {
		if !strings.Contains(out, lvl) {
			t.Errorf("missing risk level %q", lvl)
		}
	}
}

func TestFormatTable_EmptyResults(t *testing.T) {
	out := secretscore.FormatTable(nil)
	if !strings.Contains(out, "no results") {
		t.Error("expected 'no results' for empty input")
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := secretscore.FormatSummary(makeResults())
	if !strings.Contains(out, "total=3") {
		t.Errorf("expected total=3 in %q", out)
	}
	if !strings.Contains(out, "critical=1") {
		t.Errorf("expected critical=1 in %q", out)
	}
}
