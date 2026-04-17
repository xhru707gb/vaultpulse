package envelope_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/envelope"
)

func TestNew_Valid(t *testing.T) {
	e, err := envelope.New("secret/db", "v3", "enc:abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Path != "secret/db" {
		t.Errorf("expected path secret/db, got %s", e.Path)
	}
	if e.KeyVersion != "v3" {
		t.Errorf("expected key version v3, got %s", e.KeyVersion)
	}
	if e.EncryptedAt.IsZero() {
		t.Error("expected EncryptedAt to be set")
	}
}

func TestNew_MissingCiphertext(t *testing.T) {
	_, err := envelope.New("secret/db", "v3", "")
	if err == nil {
		t.Fatal("expected error for empty ciphertext")
	}
}

func TestNew_MissingKeyVersion(t *testing.T) {
	_, err := envelope.New("secret/db", "", "enc:abc")
	if err == nil {
		t.Fatal("expected error for empty key version")
	}
}

func TestAge_IsPositive(t *testing.T) {
	e, _ := envelope.New("secret/x", "v1", "enc:xyz")
	time.Sleep(2 * time.Millisecond)
	if e.Age() <= 0 {
		t.Error("expected positive age")
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	e, _ := envelope.New("secret/api", "v2", "enc:foo")
	out := envelope.FormatTable([]*envelope.Envelope{e})
	for _, hdr := range []string{"PATH", "KEY VERSION", "AGE", "ENCRYPTED AT"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatTable_ContainsPath(t *testing.T) {
	e, _ := envelope.New("secret/myapp/token", "v5", "enc:bar")
	out := envelope.FormatTable([]*envelope.Envelope{e})
	if !strings.Contains(out, "secret/myapp/token") {
		t.Error("expected path in table output")
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := envelope.FormatTable(nil)
	if !strings.Contains(out, "PATH") {
		t.Error("expected headers even for empty input")
	}
}
