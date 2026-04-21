package secretclassify_test

import (
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretclassify"
)

func defaultRules() []secretclassify.Rule {
	return []secretclassify.Rule{
		{Pattern: "prod/", Level: secretclassify.LevelSecret},
		{Pattern: "staging/", Level: secretclassify.LevelConfidential},
		{Pattern: "dev/", Level: secretclassify.LevelInternal},
	}
}

func TestNew_NoRules_ReturnsError(t *testing.T) {
	_, err := secretclassify.New(nil, secretclassify.LevelPublic)
	if err == nil {
		t.Fatal("expected error for empty rules")
	}
}

func TestNew_EmptyPattern_ReturnsError(t *testing.T) {
	_, err := secretclassify.New([]secretclassify.Rule{{Pattern: "", Level: secretclassify.LevelSecret}}, secretclassify.LevelPublic)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestClassify_MatchesRule(t *testing.T) {
	c, err := secretclassify.New(defaultRules(), secretclassify.LevelPublic)
	if err != nil {
		t.Fatal(err)
	}
	r := c.Classify("secret/prod/db")
	if r.Level != secretclassify.LevelSecret {
		t.Errorf("expected SECRET, got %s", r.Level)
	}
}

func TestClassify_FallbackApplied(t *testing.T) {
	c, _ := secretclassify.New(defaultRules(), secretclassify.LevelPublic)
	r := c.Classify("secret/other/key")
	if r.Level != secretclassify.LevelPublic {
		t.Errorf("expected PUBLIC fallback, got %s", r.Level)
	}
}

func TestClassify_FirstMatchWins(t *testing.T) {
	rules := []secretclassify.Rule{
		{Pattern: "prod/", Level: secretclassify.LevelSecret},
		{Pattern: "prod/staging", Level: secretclassify.LevelConfidential},
	}
	c, _ := secretclassify.New(rules, secretclassify.LevelPublic)
	r := c.Classify("prod/staging/key")
	if r.Level != secretclassify.LevelSecret {
		t.Errorf("expected first rule to win: SECRET, got %s", r.Level)
	}
}

func TestClassifyAll_ReturnsAllResults(t *testing.T) {
	c, _ := secretclassify.New(defaultRules(), secretclassify.LevelPublic)
	paths := []string{"secret/prod/db", "secret/dev/app", "secret/other"}
	results := c.ClassifyAll(paths)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Level != secretclassify.LevelSecret {
		t.Errorf("expected SECRET for prod path")
	}
	if results[1].Level != secretclassify.LevelInternal {
		t.Errorf("expected INTERNAL for dev path")
	}
	if results[2].Level != secretclassify.LevelPublic {
		t.Errorf("expected PUBLIC fallback")
	}
}
