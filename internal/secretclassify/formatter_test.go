package secretclassify_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpulse/internal/secretclassify"
)

func makeResults() []secretclassify.Result {
	return []secretclassify.Result{
		{Path: "secret/prod/db", Level: secretclassify.LevelSecret},
		{Path: "secret/staging/app", Level: secretclassify.LevelConfidential},
		{Path: "secret/dev/key", Level: secretclassify.LevelInternal},
		{Path: "secret/shared/token", Level: secretclassify.LevelPublic},
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := secretclassify.FormatTable(makeResults())
	if !strings.Contains(out, "PATH") {
		t.Error("expected PATH header")
	}
	if !strings.Contains(out, "LEVEL") {
		t.Error("expected LEVEL header")
	}
}

func TestFormatTable_ContainsLevelLabels(t *testing.T) {
	out := secretclassify.FormatTable(makeResults())
	for _, lbl := range []string{"SECRET", "CONFIDENTIAL", "INTERNAL", "PUBLIC"} {
		if !strings.Contains(out, lbl) {
			t.Errorf("expected label %s in output", lbl)
		}
	}
}

func TestFormatTable_EmptyResults(t *testing.T) {
	out := secretclassify.FormatTable(nil)
	if !strings.Contains(out, "no classification") {
		t.Error("expected empty message")
	}
}

func TestFormatSummary_Counts(t *testing.T) {
	out := secretclassify.FormatSummary(makeResults())
	if !strings.Contains(out, "total=4") {
		t.Errorf("expected total=4 in summary, got: %s", out)
	}
	if !strings.Contains(out, "secret=1") {
		t.Errorf("expected secret=1 in summary, got: %s", out)
	}
}
