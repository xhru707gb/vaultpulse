package redact_test

import (
	"testing"

	"github.com/yourusername/vaultpulse/internal/redact"
)

func defaultRedactor() *redact.Redactor {
	return redact.New([]string{"token", "password", "secret", "key"})
}

func TestShouldRedact_MatchesSensitiveKey(t *testing.T) {
	r := defaultRedactor()
	for _, key := range []string{"vault_token", "PASSWORD", "api_secret", "private_key"} {
		if !r.ShouldRedact(key) {
			t.Errorf("expected %q to be redacted", key)
		}
	}
}

func TestShouldRedact_AllowsSafeKey(t *testing.T) {
	r := defaultRedactor()
	for _, key := range []string{"path", "ttl", "version", "created_at"} {
		if r.ShouldRedact(key) {
			t.Errorf("expected %q NOT to be redacted", key)
		}
	}
}

func TestValue_MasksSensitive(t *testing.T) {
	r := defaultRedactor()
	got := r.Value("vault_token", "s.supersecret")
	if got != "[REDACTED]" {
		t.Fatalf("expected [REDACTED], got %q", got)
	}
}

func TestValue_PassesThroughSafe(t *testing.T) {
	r := defaultRedactor()
	got := r.Value("path", "secret/my-app")
	if got != "secret/my-app" {
		t.Fatalf("unexpected value %q", got)
	}
}

func TestMap_RedactsMatchingKeys(t *testing.T) {
	r := defaultRedactor()
	in := map[string]string{
		"path":        "secret/db",
		"db_password": "hunter2",
		"ttl":         "72h",
		"api_key":     "abc123",
	}
	out := r.Map(in)
	if out["path"] != "secret/db" {
		t.Errorf("path should be unchanged")
	}
	if out["ttl"] != "72h" {
		t.Errorf("ttl should be unchanged")
	}
	if out["db_password"] != "[REDACTED]" {
		t.Errorf("db_password should be redacted")
	}
	if out["api_key"] != "[REDACTED]" {
		t.Errorf("api_key should be redacted")
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	r := defaultRedactor()
	in := map[string]string{"token": "s.abc"}
	r.Map(in)
	if in["token"] != "s.abc" {
		t.Error("original map was mutated")
	}
}
