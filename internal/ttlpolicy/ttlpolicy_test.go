package ttlpolicy_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/ttlpolicy"
)

func defaultRules() []ttlpolicy.Rule {
	return []ttlpolicy.Rule{
		{Prefix: "secret/prod", MinTTL: 24 * time.Hour, MaxTTL: 720 * time.Hour},
		{Prefix: "secret/dev", MinTTL: time.Hour, MaxTTL: 48 * time.Hour},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := ttlpolicy.New(nil)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_MinExceedsMax_ReturnsError(t *testing.T) {
	_, err := ttlpolicy.New([]ttlpolicy.Rule{
		{Prefix: "secret/", MinTTL: 48 * time.Hour, MaxTTL: 24 * time.Hour},
	})
	if err == nil {
		t.Fatal("expected error when minTTL > maxTTL")
	}
}

func TestEvaluate_Compliant(t *testing.T) {
	e, _ := ttlpolicy.New(defaultRules())
	r := e.Evaluate("secret/prod/db", 48*time.Hour)
	if !r.Compliant {
		t.Fatalf("expected compliant, got violation: %s", r.Violation)
	}
}

func TestEvaluate_BelowMinTTL(t *testing.T) {
	e, _ := ttlpolicy.New(defaultRules())
	r := e.Evaluate("secret/prod/db", time.Hour)
	if r.Compliant {
		t.Fatal("expected violation for TTL below minimum")
	}
	if r.Violation == "" {
		t.Fatal("expected non-empty violation message")
	}
}

func TestEvaluate_ExceedsMaxTTL(t *testing.T) {
	e, _ := ttlpolicy.New(defaultRules())
	r := e.Evaluate("secret/prod/db", 1000*time.Hour)
	if r.Compliant {
		t.Fatal("expected violation for TTL above maximum")
	}
}

func TestEvaluate_NoMatchingRule_Compliant(t *testing.T) {
	e, _ := ttlpolicy.New(defaultRules())
	r := e.Evaluate("secret/staging/key", time.Minute)
	if !r.Compliant {
		t.Fatal("expected compliant when no rule matches")
	}
}

func TestEvaluateAll_ReturnsAllResults(t *testing.T) {
	e, _ := ttlpolicy.New(defaultRules())
	secrets := map[string]time.Duration{
		"secret/prod/a": 48 * time.Hour,
		"secret/dev/b":  30 * time.Minute,
	}
	results := e.EvaluateAll(secrets)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
