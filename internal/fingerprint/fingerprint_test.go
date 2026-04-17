package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/fingerprint"
)

func TestCompute_Deterministic(t *testing.T) {
	data := map[string]string{"key": "value", "foo": "bar"}
	r1, err := fingerprint.Compute("secret/a", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := fingerprint.Compute("secret/a", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Fingerprint != r2.Fingerprint {
		t.Errorf("fingerprints differ: %s vs %s", r1.Fingerprint, r2.Fingerprint)
	}
}

func TestCompute_DifferentData(t *testing.T) {
	r1, _ := fingerprint.Compute("secret/a", map[string]string{"k": "v1"})
	r2, _ := fingerprint.Compute("secret/a", map[string]string{"k": "v2"})
	if r1.Fingerprint == r2.Fingerprint {
		t.Error("expected different fingerprints for different values")
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	_, err := fingerprint.Compute("secret/a", nil)
	if err == nil {
		t.Error("expected error for nil input")
	}
}

func TestCompute_KeyCountSet(t *testing.T) {
	r, err := fingerprint.Compute("secret/a", map[string]string{"a": "1", "b": "2", "c": "3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.KeyCount != 3 {
		t.Errorf("expected KeyCount 3, got %d", r.KeyCount)
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	data := map[string]string{"secret": "old"}
	r, _ := fingerprint.Compute("secret/x", data)

	newData := map[string]string{"secret": "new"}
	changed, _, err := fingerprint.Changed("secret/x", r.Fingerprint, newData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true")
	}
}

func TestChanged_NoChange(t *testing.T) {
	data := map[string]string{"secret": "same"}
	r, _ := fingerprint.Compute("secret/x", data)
	changed, _, _ := fingerprint.Changed("secret/x", r.Fingerprint, data)
	if changed {
		t.Error("expected changed=false")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := fingerprint.FormatTable([]fingerprint.Result{
		{Path: "secret/db", Fingerprint: strings.Repeat("a", 64), KeyCount: 2},
	})
	for _, hdr := range []string{"PATH", "FINGERPRINT", "KEYS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := fingerprint.FormatTable(nil)
	if !strings.Contains(out, "No fingerprint") {
		t.Error("expected empty message")
	}
}
