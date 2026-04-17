package digest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/digest"
)

func makeReport() digest.Report {
	b := digest.NewBuilder(fixedNow)
	return b.Build([]digest.Entry{
		{Path: "secret/db", TTL: 48 * time.Hour},
		{Path: "secret/api", Expired: true},
		{Path: "secret/svc", ExpiresSoon: true, TTL: 3 * time.Hour},
		{Path: "secret/cache", Overdue: true, TTL: 1 * time.Hour},
		{Path: "secret/infra", Unhealthy: true},
	})
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	out := digest.FormatTable(makeReport())
	for _, h := range []string{"PATH", "STATE", "TTL"} {
		if !strings.Contains(out, h) {
			t.Errorf("missing header %q", h)
		}
	}
}

func TestFormatTable_StateLabels(t *testing.T) {
	out := digest.FormatTable(makeReport())
	for _, label := range []string{"EXPIRED", "EXPIRING SOON", "OVERDUE", "UNHEALTHY", "OK"} {
		if !strings.Contains(out, label) {
			t.Errorf("missing label %q", label)
		}
	}
}

func TestFormatTable_ContainsSummaryLine(t *testing.T) {
	out := digest.FormatTable(makeReport())
	if !strings.Contains(out, "Total:") {
		t.Error("missing summary line")
	}
}

func TestFormatTable_TTLFormatted(t *testing.T) {
	out := digest.FormatTable(makeReport())
	if !strings.Contains(out, "48h00m") {
		t.Errorf("expected formatted TTL in output, got:\n%s", out)
	}
}
