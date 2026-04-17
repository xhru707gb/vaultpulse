package masking_test

import (
	"testing"

	"github.com/your-org/vaultpulse/internal/masking"
)

func defaultMasker(level masking.Level) *masking.Masker {
	return masking.New(level, []string{"password", "token", "secret", "key"})
}

func TestShouldMask_SensitiveKey(t *testing.T) {
	m := defaultMasker(masking.LevelFull)
	if !m.ShouldMask("db_password") {
		t.Fatal("expected db_password to be masked")
	}
}

func TestShouldMask_SafeKey(t *testing.T) {
	m := defaultMasker(masking.LevelFull)
	if m.ShouldMask("username") {
		t.Fatal("expected username to not be masked")
	}
}

func TestMask_Full(t *testing.T) {
	m := defaultMasker(masking.LevelFull)
	got := m.Mask("api_token", "supersecret123")
	if got != "********" {
		t.Fatalf("expected ********, got %s", got)
	}
}

func TestMask_Partial(t *testing.T) {
	m := defaultMasker(masking.LevelPartial)
	got := m.Mask("api_token", "supersecret123")
	if got == "supersecret123" {
		t.Fatal("expected value to be partially masked")
	}
	if got[:2] != "su" {
		t.Fatalf("expected prefix 'su', got %s", got[:2])
	}
	if got[len(got)-2:] != "23" {
		t.Fatalf("expected suffix '23', got %s", got[len(got)-2:])
	}
}

func TestMask_None(t *testing.T) {
	m := defaultMasker(masking.LevelNone)
	got := m.Mask("api_token", "supersecret123")
	if got != "supersecret123" {
		t.Fatalf("expected plain value, got %s", got)
	}
}

func TestMask_SafeKey_Unchanged(t *testing.T) {
	m := defaultMasker(masking.LevelFull)
	got := m.Mask("username", "alice")
	if got != "alice" {
		t.Fatalf("expected alice, got %s", got)
	}
}

func TestMaskMap_MasksSensitiveKeys(t *testing.T) {
	m := defaultMasker(masking.LevelFull)
	input := map[string]string{
		"username": "alice",
		"password": "hunter2",
	}
	out := m.MaskMap(input)
	if out["username"] != "alice" {
		t.Fatalf("expected alice, got %s", out["username"])
	}
	if out["password"] != "********" {
		t.Fatalf("expected masked password, got %s", out["password"])
	}
}

func TestPartial_ShortValue(t *testing.T) {
	m := defaultMasker(masking.LevelPartial)
	got := m.Mask("secret", "ab")
	if got != "**" {
		t.Fatalf("expected **, got %s", got)
	}
}
