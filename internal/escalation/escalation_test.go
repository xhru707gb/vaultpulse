package escalation_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/escalation"
)

var defaultRules = []escalation.Rule{
	{Level: escalation.LevelCritical, Threshold: 1 * time.Hour},
	{Level: escalation.LevelWarning, Threshold: 24 * time.Hour},
	{Level: escalation.LevelInfo, Threshold: 72 * time.Hour},
}

func fixedNow() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := escalation.New(nil, fixedNow)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestEvaluate_CriticalLevel(t *testing.T) {
	e, _ := escalation.New(defaultRules, fixedNow)
	ev := e.Evaluate("secret/db", 30*time.Minute)
	if ev == nil {
		t.Fatal("expected event")
	}
	if ev.Level != escalation.LevelCritical {
		t.Errorf("expected critical, got %s", ev.Level)
	}
}

func TestEvaluate_WarningLevel(t *testing.T) {
	e, _ := escalation.New(defaultRules, fixedNow)
	ev := e.Evaluate("secret/api", 12*time.Hour)
	if ev == nil {
		t.Fatal("expected event")
	}
	if ev.Level != escalation.LevelWarning {
		t.Errorf("expected warning, got %s", ev.Level)
	}
}

func TestEvaluate_NoMatch_ReturnsNil(t *testing.T) {
	e, _ := escalation.New(defaultRules, fixedNow)
	ev := e.Evaluate("secret/safe", 168*time.Hour)
	if ev != nil {
		t.Errorf("expected nil, got %+v", ev)
	}
}

func TestEvaluateAll_ReturnsMatchingEvents(t *testing.T) {
	e, _ := escalation.New(defaultRules, fixedNow)
	secrets := map[string]time.Duration{
		"secret/a": 30 * time.Minute,
		"secret/b": 200 * time.Hour,
		"secret/c": 10 * time.Hour,
	}
	events := e.EvaluateAll(secrets)
	if len(events) != 2 {
		t.Errorf("expected 2 events, got %d", len(events))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	e, _ := escalation.New(defaultRules, fixedNow)
	events := e.EvaluateAll(map[string]time.Duration{"secret/x": 30 * time.Minute})
	out := escalation.FormatTable(events)
	for _, h := range []string{"PATH", "LEVEL", "TTL"} {
		if !contains(out, h) {
			t.Errorf("missing header %q in output", h)
		}
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := escalation.FormatTable(nil)
	if !contains(out, "No escalation") {
		t.Errorf("expected empty message, got %q", out)
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
