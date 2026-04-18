package labelmap_test

import (
	"testing"

	"github.com/yourusername/vaultpulse/internal/labelmap"
)

func TestSet_And_Get(t *testing.T) {
	m := labelmap.New()
	if err := m.Set("secret/a", "env", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels := m.Get("secret/a")
	if labels["env"] != "prod" {
		t.Errorf("expected prod, got %s", labels["env"])
	}
}

func TestSet_EmptyKey_ReturnsError(t *testing.T) {
	m := labelmap.New()
	if err := m.Set("secret/a", "", "val"); err == nil {
		t.Error("expected error for empty key")
	}
}

func TestGet_UnknownPath_ReturnsNil(t *testing.T) {
	m := labelmap.New()
	if got := m.Get("secret/missing"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	m := labelmap.New()
	_ = m.Set("secret/a", "env", "prod")
	m.Delete("secret/a", "env")
	if got := m.Get("secret/a"); got != nil {
		t.Errorf("expected nil after delete, got %v", got)
	}
}

func TestFilter_MatchesSelector(t *testing.T) {
	m := labelmap.New()
	_ = m.Set("secret/a", "env", "prod")
	_ = m.Set("secret/a", "team", "ops")
	_ = m.Set("secret/b", "env", "dev")
	_ = m.Set("secret/c", "env", "prod")

	results := m.Filter(map[string]string{"env": "prod"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0] != "secret/a" || results[1] != "secret/c" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	m := labelmap.New()
	_ = m.Set("secret/a", "env", "dev")
	results := m.Filter(map[string]string{"env": "prod"})
	if len(results) != 0 {
		t.Errorf("expected empty, got %v", results)
	}
}

func TestFormatLabels_SortedOutput(t *testing.T) {
	labels := map[string]string{"team": "ops", "env": "prod"}
	got := labelmap.FormatLabels(labels)
	expected := "env=prod, team=ops"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFormatLabels_Empty(t *testing.T) {
	got := labelmap.FormatLabels(map[string]string{})
	if got != "<none>" {
		t.Errorf("expected <none>, got %q", got)
	}
}
