package secretbundle_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretbundle"
)

func newRegistry(t *testing.T) *secretbundle.Registry {
	t.Helper()
	return secretbundle.New()
}

func TestAdd_And_Evaluate(t *testing.T) {
	r := newRegistry(t)
	if err := r.Add("infra"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = r.AddEntry("infra", "secret/db", 1, false)
	_ = r.AddEntry("infra", "secret/api", 2, false)
	res, err := r.Evaluate("infra")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Total != 2 || res.Expired != 0 || !res.Healthy {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	r := newRegistry(t)
	if err := r.Add(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_Duplicate_ReturnsError(t *testing.T) {
	r := newRegistry(t)
	_ = r.Add("svc")
	if err := r.Add("svc"); err == nil {
		t.Fatal("expected error for duplicate bundle")
	}
}

func TestAddEntry_UnknownBundle_ReturnsError(t *testing.T) {
	r := newRegistry(t)
	if err := r.AddEntry("missing", "secret/x", 1, false); err == nil {
		t.Fatal("expected error for unknown bundle")
	}
}

func TestAddEntry_EmptyPath_ReturnsError(t *testing.T) {
	r := newRegistry(t)
	_ = r.Add("b")
	if err := r.AddEntry("b", "", 1, false); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestEvaluate_ExpiredCounted(t *testing.T) {
	r := newRegistry(t)
	_ = r.Add("creds")
	_ = r.AddEntry("creds", "secret/a", 1, false)
	_ = r.AddEntry("creds", "secret/b", 1, true)
	res, _ := r.Evaluate("creds")
	if res.Expired != 1 || res.Healthy {
		t.Errorf("expected 1 expired and unhealthy, got %+v", res)
	}
}

func TestEvaluateAll_ReturnsAllBundles(t *testing.T) {
	r := newRegistry(t)
	_ = r.Add("a")
	_ = r.Add("b")
	results := r.EvaluateAll()
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretbundle.FormatTable([]secretbundle.EvalResult{
		{Name: "infra", Total: 3, Expired: 1, Healthy: false},
	})
	for _, hdr := range []string{"BUNDLE", "TOTAL", "EXPIRED", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	results := []secretbundle.EvalResult{
		{Name: "a", Healthy: true},
		{Name: "b", Healthy: false},
	}
	out := secretbundle.FormatSummary(results)
	if !strings.Contains(out, "2 total") {
		t.Errorf("expected total count in summary, got: %s", out)
	}
	if !strings.Contains(out, "1 degraded") {
		t.Errorf("expected degraded count in summary, got: %s", out)
	}
}
