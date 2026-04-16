package notify_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/notify"
)

// stubSender records calls and optionally returns an error.
type stubSender struct {
	name    notify.Channel
	calls   []notify.Message
	failErr error
}

func (s *stubSender) Name() notify.Channel { return s.name }
func (s *stubSender) Send(_ context.Context, msg notify.Message) error {
	s.calls = append(s.calls, msg)
	return s.failErr
}

func TestDispatch_CallsAllSenders(t *testing.T) {
	a := &stubSender{name: "a"}
	b := &stubSender{name: "b"}
	d := notify.New(a, b)

	msg := notify.Message{Level: "warning", Subject: "secret/db", Body: "expires soon"}
	if err := d.Dispatch(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.calls) != 1 || len(b.calls) != 1 {
		t.Fatalf("expected 1 call each, got a=%d b=%d", len(a.calls), len(b.calls))
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	fail := &stubSender{name: "fail", failErr: errors.New("boom")}
	ok := &stubSender{name: "ok"}
	d := notify.New(fail, ok)

	err := d.Dispatch(context.Background(), notify.Message{Level: "critical"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// ok sender should still have been called
	if len(ok.calls) != 1 {
		t.Fatalf("expected ok sender to be called, got %d calls", len(ok.calls))
	}
}

func TestStdoutSender_WritesLevel(t *testing.T) {
	var buf bytes.Buffer
	s := notify.NewStdoutSender(&buf)

	msg := notify.Message{Level: "critical", Subject: "secret/api", Body: "expired"}
	if err := s.Send(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CRITICAL") {
		t.Errorf("expected CRITICAL in output, got: %s", out)
	}
	if !strings.Contains(out, "secret/api") {
		t.Errorf("expected subject in output, got: %s", out)
	}
}

func TestRegisterSender_AddsAtRuntime(t *testing.T) {
	d := notify.New()
	s := &stubSender{name: "late"}
	d.RegisterSender(s)

	_ = d.Dispatch(context.Background(), notify.Message{Level: "info", Subject: "x"})
	if len(s.calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(s.calls))
	}
}
