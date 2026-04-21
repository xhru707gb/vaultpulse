package secretrelay_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretrelay"
)

func newRelay() *secretrelay.Relay { return secretrelay.New() }

func TestRegister_And_Len(t *testing.T) {
	r := newRelay()
	if err := r.Register("sink1", func(_ string, _ map[string]string) error { return nil }); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := r.Len(); got != 1 {
		t.Fatalf("expected 1 handler, got %d", got)
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	r := newRelay()
	err := r.Register("", func(_ string, _ map[string]string) error { return nil })
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_NilHandler_ReturnsError(t *testing.T) {
	r := newRelay()
	if err := r.Register("h", nil); err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	r := newRelay()
	h := func(_ string, _ map[string]string) error { return nil }
	_ = r.Register("dup", h)
	if err := r.Register("dup", h); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestDeregister_RemovesHandler(t *testing.T) {
	r := newRelay()
	_ = r.Register("h", func(_ string, _ map[string]string) error { return nil })
	r.Deregister("h")
	if r.Len() != 0 {
		t.Fatal("expected 0 handlers after deregister")
	}
}

func TestDeregister_Unknown_NoOp(t *testing.T) {
	r := newRelay()
	r.Deregister("nonexistent") // should not panic
}

func TestDispatch_CallsAllHandlers(t *testing.T) {
	r := newRelay()
	called := map[string]bool{}
	for _, name := range []string{"a", "b", "c"} {
		n := name
		_ = r.Register(n, func(_ string, _ map[string]string) error {
			called[n] = true
			return nil
		})
	}
	errs := r.Dispatch("secret/path", map[string]string{"key": "val"})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(called) != 3 {
		t.Fatalf("expected 3 handlers called, got %d", len(called))
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	r := newRelay()
	_ = r.Register("ok", func(_ string, _ map[string]string) error { return nil })
	_ = r.Register("bad", func(_ string, _ map[string]string) error { return errors.New("boom") })
	errs := r.Dispatch("secret/path", nil)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestDispatch_PayloadForwarded(t *testing.T) {
	r := newRelay()
	var got map[string]string
	_ = r.Register("capture", func(_ string, p map[string]string) error {
		got = p
		return nil
	})
	want := map[string]string{"token": "abc123"}
	r.Dispatch("secret/token", want)
	if got["token"] != "abc123" {
		t.Fatalf("payload not forwarded correctly: %v", got)
	}
}
