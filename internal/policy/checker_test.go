package policy

import (
	"testing"
	"time"
)

var fixedNow = func() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func defaultPolicies() []Policy {
	return []Policy{
		{
			Name:           "standard",
			MaxTTLDays:     30,
			RequireRotation: true,
			RotationDays:   90,
		},
	}
}

func TestEvaluate_Compliant(t *testing.T) {
	c := NewChecker(defaultPolicies(), fixedNow)
	statuses := c.Evaluate("secret/db", 10*24*time.Hour, fixedNow().Add(-30*24*time.Hour))
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if !statuses[0].Compliant {
		t.Errorf("expected compliant, got violations: %v", statuses[0].Violations)
	}
}

func TestEvaluate_TTLViolation(t *testing.T) {
	c := NewChecker(defaultPolicies(), fixedNow)
	statuses := c.Evaluate("secret/db", 60*24*time.Hour, fixedNow().Add(-10*24*time.Hour))
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].Compliant {
		t.Error("expected non-compliant for TTL violation")
	}
	if len(statuses[0].Violations) == 0 {
		t.Error("expected at least one violation message")
	}
}

func TestEvaluate_RotationViolation(t *testing.T) {
	c := NewChecker(defaultPolicies(), fixedNow)
	oldRotation := fixedNow().Add(-120 * 24 * time.Hour)
	statuses := c.Evaluate("secret/api", 5*24*time.Hour, oldRotation)
	if statuses[0].Compliant {
		t.Error("expected non-compliant for rotation violation")
	}
	found := false
	for _, v := range statuses[0].Violations {
		if len(v) > 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected rotation violation message")
	}
}

func TestEvaluate_NoLastRotated_SkipsRotationCheck(t *testing.T) {
	c := NewChecker(defaultPolicies(), fixedNow)
	statuses := c.Evaluate("secret/new", 5*24*time.Hour, time.Time{})
	if !statuses[0].Compliant {
		t.Errorf("expected compliant when lastRotated is zero, got: %v", statuses[0].Violations)
	}
}

func TestEvaluateAll_MultipleSecrets(t *testing.T) {
	c := NewChecker(defaultPolicies(), fixedNow)
	secrets := map[string]SecretMeta{
		"secret/a": {TTL: 5 * 24 * time.Hour, LastRotated: fixedNow().Add(-10 * 24 * time.Hour)},
		"secret/b": {TTL: 60 * 24 * time.Hour, LastRotated: fixedNow().Add(-10 * 24 * time.Hour)},
	}
	all := c.EvaluateAll(secrets)
	if len(all) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(all))
	}
}
