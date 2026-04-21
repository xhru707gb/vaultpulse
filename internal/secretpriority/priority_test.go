package secretpriority_test

import (
	"strings"
	"testing"

	"vaultpulse/internal/secretpriority"
)

func defaultRules() []secretpriority.Rule {
	return []secretpriority.Rule{
		{Prefix: "secret/prod/", Level: secretpriority.LevelCritical},
		{Prefix: "secret/staging/", Level: secretpriority.LevelHigh},
		{Prefix: "secret/dev/", Level: secretpriority.LevelMedium},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secretpriority.New(nil, secretpriority.LevelLow)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestEvaluate_MatchesCritical(t *testing.T) {
	e, err := secretpriority.New(defaultRules(), secretpriority.LevelLow)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := e.Evaluate("secret/prod/db-password")
	if r.Level != secretpriority.LevelCritical {
		t.Errorf("expected CRITICAL, got %v", r.Level)
	}
	if r.Rule != "secret/prod/" {
		t.Errorf("expected rule secret/prod/, got %q", r.Rule)
	}
}

func TestEvaluate_FallsBackToDefault(t *testing.T) {
	e, _ := secretpriority.New(defaultRules(), secretpriority.LevelLow)
	r := e.Evaluate("secret/other/token")
	if r.Level != secretpriority.LevelLow {
		t.Errorf("expected LOW, got %v", r.Level)
	}
	if r.Rule != "default" {
		t.Errorf("expected rule 'default', got %q", r.Rule)
	}
}

func TestEvaluateAll_ReturnsAllResults(t *testing.T) {
	e, _ := secretpriority.New(defaultRules(), secretpriority.LevelLow)
	paths := []string{"secret/prod/key", "secret/dev/key", "secret/unknown/key"}
	results := e.EvaluateAll(paths)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Level != secretpriority.LevelCritical {
		t.Errorf("first result should be CRITICAL")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	e, _ := secretpriority.New(defaultRules(), secretpriority.LevelLow)
	results := e.EvaluateAll([]string{"secret/prod/api-key"})
	out := secretpriority.FormatTable(results)
	for _, hdr := range []string{"PATH", "PRIORITY", "MATCHED RULE"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_EmptyResults(t *testing.T) {
	out := secretpriority.FormatTable(nil)
	if !strings.Contains(out, "No priority results") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	e, _ := secretpriority.New(defaultRules(), secretpriority.LevelLow)
	results := e.EvaluateAll([]string{
		"secret/prod/a", "secret/prod/b",
		"secret/staging/c",
		"secret/dev/d",
		"secret/other/e",
	})
	summary := secretpriority.FormatSummary(results)
	if !strings.Contains(summary, "Total: 5") {
		t.Errorf("expected Total: 5, got: %s", summary)
	}
	if !strings.Contains(summary, "Critical: 2") {
		t.Errorf("expected Critical: 2, got: %s", summary)
	}
}
