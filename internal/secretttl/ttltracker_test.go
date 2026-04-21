package secretttl_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpulse/internal/secretttl"
)

var fixedNow = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func newTracker() *secretttl.Tracker {
	return secretttl.New(func() time.Time { return fixedNow })
}

func baseEntry(path string, expiresIn, warnIn time.Duration) secretttl.Entry {
	return secretttl.Entry{
		Path:      path,
		ExpiresAt: fixedNow.Add(expiresIn),
		WarningIn: warnIn,
	}
}

func TestRegister_And_Evaluate_OK(t *testing.T) {
	tr := newTracker()
	err := tr.Register(baseEntry("secret/a", 48*time.Hour, 24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, ok := tr.Evaluate("secret/a")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if s.State != secretttl.StateOK {
		t.Errorf("expected OK, got %v", s.State)
	}
}

func TestEvaluate_Warning(t *testing.T) {
	tr := newTracker()
	_ = tr.Register(baseEntry("secret/b", 6*time.Hour, 24*time.Hour))
	s, _ := tr.Evaluate("secret/b")
	if s.State != secretttl.StateWarning {
		t.Errorf("expected Warning, got %v", s.State)
	}
}

func TestEvaluate_Expired(t *testing.T) {
	tr := newTracker()
	_ = tr.Register(baseEntry("secret/c", -1*time.Second, 24*time.Hour))
	s, _ := tr.Evaluate("secret/c")
	if s.State != secretttl.StateExpired {
		t.Errorf("expected Expired, got %v", s.State)
	}
}

func TestEvaluate_NotFound(t *testing.T) {
	tr := newTracker()
	_, ok := tr.Evaluate("secret/missing")
	if ok {
		t.Error("expected false for unknown path")
	}
}

func TestRegister_EmptyPath_ReturnsError(t *testing.T) {
	tr := newTracker()
	err := tr.Register(secretttl.Entry{Path: "", ExpiresAt: fixedNow.Add(time.Hour), WarningIn: time.Minute})
	if err == nil {
		t.Error("expected error for empty path")
	}
}

func TestRegister_ZeroExpiry_ReturnsError(t *testing.T) {
	tr := newTracker()
	err := tr.Register(secretttl.Entry{Path: "secret/x", WarningIn: time.Minute})
	if err == nil {
		t.Error("expected error for zero expiresAt")
	}
}

func TestRegister_InvalidWarningIn_ReturnsError(t *testing.T) {
	tr := newTracker()
	err := tr.Register(secretttl.Entry{Path: "secret/x", ExpiresAt: fixedNow.Add(time.Hour), WarningIn: 0})
	if err == nil {
		t.Error("expected error for zero warningIn")
	}
}

func TestEvaluateAll_ReturnsBoth(t *testing.T) {
	tr := newTracker()
	_ = tr.Register(baseEntry("secret/ok", 48*time.Hour, 24*time.Hour))
	_ = tr.Register(baseEntry("secret/exp", -time.Second, 24*time.Hour))
	all := tr.EvaluateAll()
	if len(all) != 2 {
		t.Errorf("expected 2 statuses, got %d", len(all))
	}
}
