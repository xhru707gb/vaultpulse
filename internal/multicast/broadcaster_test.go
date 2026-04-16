package multicast_test

import (
	"sync/atomic"
	"testing"

	"github.com/your-org/vaultpulse/internal/multicast"
)

func TestRegister_IncreasesLen(t *testing.T) {
	b := multicast.New()
	b.Register("a", func(string, any) {})
	b.Register("b", func(string, any) {})
	if got := b.Len(); got != 2 {
		t.Fatalf("expected 2 handlers, got %d", got)
	}
}

func TestRegister_NilOrEmptyIgnored(t *testing.T) {
	b := multicast.New()
	b.Register("", func(string, any) {})
	b.Register("a", nil)
	if b.Len() != 0 {
		t.Fatal("expected no handlers registered")
	}
}

func TestDeregister_RemovesHandler(t *testing.T) {
	b := multicast.New()
	b.Register("x", func(string, any) {})
	b.Deregister("x")
	if b.Len() != 0 {
		t.Fatal("expected handler to be removed")
	}
}

func TestDeregister_UnknownName_NoOp(t *testing.T) {
	b := multicast.New()
	b.Deregister("ghost") // must not panic
}

func TestBroadcast_CallsAllHandlers(t *testing.T) {
	b := multicast.New()
	var countA, countB atomic.Int32
	b.Register("a", func(e string, _ any) { countA.Add(1) })
	b.Register("b", func(e string, _ any) { countB.Add(1) })

	b.Broadcast("expiry.warning", "payload")

	if countA.Load() != 1 || countB.Load() != 1 {
		t.Fatalf("expected each handler called once, got a=%d b=%d", countA.Load(), countB.Load())
	}
}

func TestBroadcast_PassesEventAndPayload(t *testing.T) {
	b := multicast.New()
	var gotEvent string
	var gotPayload any
	b.Register("h", func(e string, p any) {
		gotEvent = e
		gotPayload = p
	})
	b.Broadcast("health.critical", 42)
	if gotEvent != "health.critical" {
		t.Fatalf("unexpected event: %s", gotEvent)
	}
	if gotPayload.(int) != 42 {
		t.Fatalf("unexpected payload: %v", gotPayload)
	}
}

func TestBroadcast_NoHandlers_NoOp(t *testing.T) {
	b := multicast.New()
	b.Broadcast("any", nil) // must not panic or block
}
