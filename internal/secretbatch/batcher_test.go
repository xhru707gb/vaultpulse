package secretbatch_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/vaultpulse/vaultpulse/internal/secretbatch"
)

func okEvaluator(path string) (bool, string, error) {
	return true, "all good", nil
}

func failEvaluator(path string) (bool, string, error) {
	return false, "check failed", nil
}

func errEvaluator(path string) (bool, string, error) {
	return false, "", errors.New("vault unreachable")
}

func TestNew_InvalidConcurrency(t *testing.T) {
	_, err := secretbatch.New(0, okEvaluator)
	if err == nil {
		t.Fatal("expected error for concurrency=0")
	}
}

func TestNew_NilEvaluator(t *testing.T) {
	_, err := secretbatch.New(2, nil)
	if err == nil {
		t.Fatal("expected error for nil evaluator")
	}
}

func TestRun_AllOK(t *testing.T) {
	b, err := secretbatch.New(2, okEvaluator)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	paths := []string{"secret/a", "secret/b", "secret/c"}
	results := b.Run(paths)

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.OK {
			t.Errorf("expected OK for path %s", r.Path)
		}
	}
}

func TestRun_PreservesOrder(t *testing.T) {
	b, _ := secretbatch.New(4, okEvaluator)
	paths := []string{"a", "b", "c", "d", "e"}
	results := b.Run(paths)

	for i, r := range results {
		if r.Path != paths[i] {
			t.Errorf("index %d: expected path %s, got %s", i, paths[i], r.Path)
		}
	}
}

func TestRun_EmptyPaths(t *testing.T) {
	b, _ := secretbatch.New(1, okEvaluator)
	results := b.Run(nil)
	if results != nil {
		t.Errorf("expected nil results for empty input")
	}
}

func TestFailures_FiltersCorrectly(t *testing.T) {
	b, _ := secretbatch.New(2, failEvaluator)
	results := b.Run([]string{"secret/x", "secret/y"})
	fails := secretbatch.Failures(results)
	if len(fails) != 2 {
		t.Errorf("expected 2 failures, got %d", len(fails))
	}
}

func TestSummary_ContainsCounts(t *testing.T) {
	b, _ := secretbatch.New(1, errEvaluator)
	results := b.Run([]string{"secret/a", "secret/b"})
	s := secretbatch.Summary(results)
	if !strings.Contains(s, "2 evaluated") {
		t.Errorf("expected summary to mention total, got: %s", s)
	}
	if !strings.Contains(s, "2 failed") {
		t.Errorf("expected summary to mention failures, got: %s", s)
	}
}
