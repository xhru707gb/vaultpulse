package circuit_test

import (
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/circuit"
)

func fixedNow(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestNew_InvalidConfig(t *testing.T) {
	_, err := circuit.New(circuit.Config{MaxFailures: 0, OpenTimeout: time.Second})
	if err == nil {
		t.Fatal("expected error for MaxFailures=0")
	}
	_, err = circuit.New(circuit.Config{MaxFailures: 1, OpenTimeout: 0})
	if err == nil {
		t.Fatal("expected error for OpenTimeout=0")
	}
}

func TestAllow_ClosedByDefault(t *testing.T) {
	br, _ := circuit.New(circuit.Config{MaxFailures: 3, OpenTimeout: time.Second})
	if err := br.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensCircuit(t *testing.T) {
	br, _ := circuit.New(circuit.Config{MaxFailures: 2, OpenTimeout: time.Second})
	br.RecordFailure()
	if br.CurrentState() != circuit.StateClosed {
		t.Fatal("expected closed after 1 failure")
	}
	br.RecordFailure()
	if br.CurrentState() != circuit.StateOpen {
		t.Fatal("expected open after 2 failures")
	}
	if err := br.Allow(); err != circuit.ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestRecordSuccess_ResetsBreakerFromHalfOpen(t *testing.T) {
	now := time.Now()
	br, _ := circuit.New(circuit.Config{MaxFailures: 1, OpenTimeout: 5 * time.Second})
	br.SetNow(fixedNow(now))
	br.RecordFailure()
	// advance past open timeout
	br.SetNow(fixedNow(now.Add(6 * time.Second)))
	if err := br.Allow(); err != nil {
		t.Fatalf("expected half-open allow, got %v", err)
	}
	br.RecordSuccess()
	if br.CurrentState() != circuit.StateClosed {
		t.Fatal("expected closed after success")
	}
	if br.Failures() != 0 {
		t.Fatal("expected failures reset to 0")
	}
}

func TestAllow_StillOpenBeforeTimeout(t *testing.T) {
	now := time.Now()
	br, _ := circuit.New(circuit.Config{MaxFailures: 1, OpenTimeout: 10 * time.Second})
	br.SetNow(fixedNow(now))
	br.RecordFailure()
	br.SetNow(fixedNow(now.Add(3 * time.Second)))
	if err := br.Allow(); err != circuit.ErrOpen {
		t.Fatalf("expected ErrOpen before timeout, got %v", err)
	}
}
