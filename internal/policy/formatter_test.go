package policy

import (
	"strings"
	"testing"
	"time"
)

func makeStatus(path string, compliant bool, violations ...string) Status {
	return Status{
		Path:       path,
		Policy:     "standard",
		Compliant:  compliant,
		Violations: violations,
		CheckedAt:  time.Now(),
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := FormatTable([]Status{makeStatus("secret/db", true)})
	for _, h := range []string{"PATH", "POLICY", "STATUS", "VIOLATIONS"} {
		if !strings.Contains(out, h) {
			t.Errorf("expected header %q in output", h)
		}
	}
}

func TestFormatTable_ComplianceLabels(t *testing.T) {
	statuses := []Status{
		makeStatus("secret/ok", true),
		makeStatus("secret/bad", false, "TTL too long"),
	}
	out := FormatTable(statuses)
	if !strings.Contains(out, labelCompliant) {
		t.Errorf("expected %q in output", labelCompliant)
	}
	if !strings.Contains(out, labelViolation) {
		t.Errorf("expected %q in output", labelViolation)
	}
}

func TestFormatTable_EmptyStatuses(t *testing.T) {
	out := FormatTable(nil)
	if !strings.Contains(out, "No policy") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatTable_ViolationMessage(t *testing.T) {
	s := makeStatus("secret/api", false, "TTL 720h exceeds max 30d", "last rotated 120d ago")
	out := FormatTable([]Status{s})
	if !strings.Contains(out, "TTL 720h exceeds max 30d") {
		t.Errorf("expected violation message in output, got:\n%s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	statuses := []Status{
		makeStatus("secret/a", true),
		makeStatus("secret/b", false, "bad"),
		makeStatus("secret/c", false, "worse"),
	}
	out := FormatSummary(statuses)
	if !strings.Contains(out, "3 evaluated") {
		t.Errorf("expected total count, got: %s", out)
	}
	if !strings.Contains(out, "1 compliant") {
		t.Errorf("expected compliant count, got: %s", out)
	}
	if !strings.Contains(out, "2 violations") {
		t.Errorf("expected violation count, got: %s", out)
	}
}
