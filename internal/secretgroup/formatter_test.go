package secretgroup_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretgroup"
)

func buildGroups() []*secretgroup.Group {
	g := secretgroup.New()
	_ = g.Add("infra", "secret/db/prod")
	_ = g.Add("infra", "secret/db/staging")
	_ = g.Add("app", "secret/app/config")
	return g.All()
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretgroup.FormatTable(buildGroups())
	for _, h := range []string{"GROUP", "PATHS", "SAMPLE PATH"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q in output:\n%s", h, out)
		}
	}
}

func TestFormatTable_ContainsGroupNames(t *testing.T) {
	out := secretgroup.FormatTable(buildGroups())
	if !strings.Contains(out, "infra") || !strings.Contains(out, "app") {
		t.Errorf("missing group names in output:\n%s", out)
	}
}

func TestFormatTable_EmptyGroups(t *testing.T) {
	out := secretgroup.FormatTable(nil)
	if !strings.Contains(out, "no groups") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := secretgroup.FormatSummary(buildGroups())
	if !strings.Contains(out, "2 group") {
		t.Errorf("expected 2 groups in summary, got: %s", out)
	}
	if !strings.Contains(out, "3 total") {
		t.Errorf("expected 3 total paths in summary, got: %s", out)
	}
}
