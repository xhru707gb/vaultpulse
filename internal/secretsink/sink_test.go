package secretsink_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretsink"
)

// --- helpers ---

type mockSink struct {
	name   string
	events []secretsink.Event
	Err    error
}

func (m *mockSink) Name() string { return m.name }
func (m *mockSink) Send(e secretsink.Event) error {
	if m.Err != nil {
		return m.Err
	}
	m.events = append(m.events, e)
	return nil
}

// --- tests ---

func TestRegister_And_Len(t *testing.T) {
	r := secretsink.New()
	if err := r.Register(&mockSink{name: "stdout"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 sink, got %d", r.Len())
	}
}

func TestRegister_Duplicate_ReturnsError(t *testing.T) {
	r := secretsink.New()
	s := &mockSink{name: "file"}
	_ = r.Register(s)
	if err := r.Register(s); err == nil {
		t.Fatal("expected error for duplicate sink name")
	}
}

func TestRegister_NilSink_ReturnsError(t *testing.T) {
	r := secretsink.New()
	if err := r.Register(nil); err == nil {
		t.Fatal("expected error for nil sink")
	}
}

func TestDeregister_RemovesSink(t *testing.T) {
	r := secretsink.New()
	_ = r.Register(&mockSink{name: "s1"})
	r.Deregister("s1")
	if r.Len() != 0 {
		t.Fatalf("expected 0 sinks after deregister, got %d", r.Len())
	}
}

func TestDispatch_CallsAllSinks(t *testing.T) {
	r := secretsink.New()
	a := &mockSink{name: "a"}
	b := &mockSink{name: "b"}
	_ = r.Register(a)
	_ = r.Register(b)

	e := secretsink.Event{Path: "secret/db", Kind: "expired", Message: "TTL reached zero"}
	errs := r.Dispatch(e)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(a.events) != 1 || len(b.events) != 1 {
		t.Fatal("expected each sink to receive exactly one event")
	}
}

func TestDispatch_CollectsSinkErrors(t *testing.T) {
	r := secretsink.New()
	bad := &mockSink{name: "bad", Err: errors.New("send failed")}
	_ = r.Register(bad)

	errs := r.Dispatch(secretsink.Event{Path: "p", Kind: "warning"})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}
