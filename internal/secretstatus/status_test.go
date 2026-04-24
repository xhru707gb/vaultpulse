package secretstatus_test

import (
	"errors"
	"testing"

	"vaultpulse/internal/secretstatus"
)

// stubProvider is a test-only Provider implementation.
type stubProvider struct {
	name   string
	level  secretstatus.Level
	reason string
	err    error
}

func (s *stubProvider) Name() string { return s.name }
func (s *stubProvider) Evaluate(_ string) (secretstatus.Level, string, error) {
	return s.level, s.reason, s.err
}

func TestNew_NoProviders_ReturnsError(t *testing.T) {
	_, err := secretstatus.New()
	if err == nil {
		t.Fatal("expected error for zero providers")
	}
}

func TestEvaluate_AllOK(t *testing.T) {
	p := &stubProvider{name: "expiry", level: secretstatus.LevelOK, reason: ""}
	eval, _ := secretstatus.New(p)
	entry, err := eval.Evaluate("secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != secretstatus.LevelOK {
		t.Errorf("expected OK, got %v", entry.Level)
	}
	if len(entry.Reasons) != 0 {
		t.Errorf("expected no reasons, got %v", entry.Reasons)
	}
}

func TestEvaluate_TakesHighestLevel(t *testing.T) {
	p1 := &stubProvider{name: "expiry", level: secretstatus.LevelWarning, reason: "expires soon"}
	p2 := &stubProvider{name: "rotation", level: secretstatus.LevelCritical, reason: "overdue"}
	eval, _ := secretstatus.New(p1, p2)
	entry, err := eval.Evaluate("secret/bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != secretstatus.LevelCritical {
		t.Errorf("expected CRITICAL, got %v", entry.Level)
	}
	if len(entry.Reasons) != 2 {
		t.Errorf("expected 2 reasons, got %d", len(entry.Reasons))
	}
}

func TestEvaluate_ProviderErrorSkipped(t *testing.T) {
	p := &stubProvider{name: "policy", err: errors.New("unavailable")}
	eval, _ := secretstatus.New(p)
	entry, err := eval.Evaluate("secret/baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != secretstatus.LevelOK {
		t.Errorf("expected OK when provider errors, got %v", entry.Level)
	}
}

func TestEvaluate_EmptyPath_ReturnsError(t *testing.T) {
	p := &stubProvider{name: "expiry", level: secretstatus.LevelOK}
	eval, _ := secretstatus.New(p)
	_, err := eval.Evaluate("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestEvaluateAll_ReturnsAllEntries(t *testing.T) {
	p := &stubProvider{name: "expiry", level: secretstatus.LevelWarning, reason: "soon"}
	eval, _ := secretstatus.New(p)
	paths := []string{"secret/a", "secret/b", "secret/c"}
	entries, err := eval.EvaluateAll(paths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestEvaluateAll_EmptyPaths_ReturnsError(t *testing.T) {
	p := &stubProvider{name: "expiry", level: secretstatus.LevelOK}
	eval, _ := secretstatus.New(p)
	_, err := eval.EvaluateAll(nil)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}
