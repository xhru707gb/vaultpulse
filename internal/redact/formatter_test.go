package redact_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/redact"
)

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := redact.FormatTable(map[string]string{"path": "secret/app"})
	if !strings.Contains(out, "FIELD") {
		t.Error("expected FIELD header")
	}
	if !strings.Contains(out, "VALUE") {
		t.Error("expected VALUE header")
	}
}

func TestFormatTable_ContainsEntry(t *testing.T) {
	out := redact.FormatTable(map[string]string{"ttl": "48h"})
	if !strings.Contains(out, "ttl") {
		t.Error("expected key in output")
	}
	if !strings.Contains(out, "48h") {
		t.Error("expected value in output")
	}
}

func TestFormatTable_EmptyMap(t *testing.T) {
	out := redact.FormatTable(map[string]string{})
	if !strings.Contains(out, "no fields") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatTable_RedactedValue(t *testing.T) {
	r := defaultRedactor()
	m := r.Map(map[string]string{"vault_token": "s.secret", "path": "kv/app"})
	out := redact.FormatTable(m)
	if strings.Contains(out, "s.secret") {
		t.Error("raw secret value should not appear in output")
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Error("expected [REDACTED] in output")
	}
}
