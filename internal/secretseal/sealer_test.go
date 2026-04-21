package secretseal_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretseal"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTestSealer() *secretseal.Sealer {
	return secretseal.New()
}

func TestSeal_And_IsSealed(t *testing.T) {
	s := newTestSealer()
	if err := s.Seal("secret/db", "compliance", fixedNow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsSealed("secret/db") {
		t.Error("expected secret/db to be sealed")
	}
}

func TestSeal_EmptyPath_ReturnsError(t *testing.T) {
	s := newTestSealer()
	if err := s.Seal("", "reason", fixedNow); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestSeal_DuplicatePath_ReturnsError(t *testing.T) {
	s := newTestSealer()
	if err := s.Seal("secret/db", "compliance", fixedNow); err != nil {
		t.Fatalf("unexpected error on first seal: %v", err)
	}
	if err := s.Seal("secret/db", "duplicate", fixedNow); err == nil {
		t.Error("expected error when sealing an already-sealed path")
	}
}

func TestUnseal_RemovesSeal(t *testing.T) {
	s := newTestSealer()
	_ = s.Seal("secret/api", "audit", fixedNow)
	if err := s.Unseal("secret/api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsSealed("secret/api") {
		t.Error("expected secret/api to be unsealed")
	}
}

func TestUnseal_UnknownPath_ReturnsError(t *testing.T) {
	s := newTestSealer()
	if err := s.Unseal("secret/missing"); err == nil {
		t.Error("expected error for unknown path")
	}
}

func TestEvaluate_ReturnsAllStatuses(t *testing.T) {
	s := newTestSealer()
	_ = s.Seal("secret/a", "reason-a", fixedNow)
	_ = s.Seal("secret/b", "reason-b", fixedNow)

	statuses := s.Evaluate()
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}
	for _, st := range statuses {
		if !st.Sealed {
			t.Errorf("expected %q to be sealed", st.Path)
		}
		if st.SealedAt.IsZero() {
			t.Errorf("expected SealedAt to be set for %q", st.Path)
		}
	}
}

func TestEvaluate_ReflectsUnsealedState(t *testing.T) {
	s := newTestSealer()
	_ = s.Seal("secret/x", "test", fixedNow)
	_ = s.Unseal("secret/x")

	statuses := s.Evaluate()
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].Sealed {
		t.Error("expected status to show unsealed")
	}
}

func TestIsSealed_UnknownPath_ReturnsFalse(t *testing.T) {
	s := newTestSealer()
	if s.IsSealed("secret/nope") {
		t.Error("expected false for unknown path")
	}
}
