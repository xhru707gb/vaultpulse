package digest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpulse/internal/digest"
)

func fixedNow() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func makeEntries() []digest.Entry {
	return []digest.Entry{
		{Path: "secret/db", Expired: false, ExpiresSoon: false, TTL: 48 * time.Hour},
		{Path: "secret/api", Expired: true, TTL: 0},
		{Path: "secret/cache", Overdue: true, TTL: 2 * time.Hour},
	}
}

func TestBuild_AlertCount(t *testing.T) {
	b := digest.NewBuilder(fixedNow)
	r := b.Build(makeEntries())
	if r.AlertCount != 2 {
		t.Fatalf("expected 2 alerts, got %d", r.AlertCount)
	}
}

func TestBuild_TotalSecrets(t *testing.T) {
	b := digest.NewBuilder(fixedNow)
	r := b.Build(makeEntries())
	if r.TotalSecrets != 3 {
		t.Fatalf("expected 3 secrets, got %d", r.TotalSecrets)
	}
}

func TestBuild_GeneratedAt(t *testing.T) {
	b := digest.NewBuilder(fixedNow)
	r := b.Build(makeEntries())
	if !r.GeneratedAt.Equal(fixedNow()) {
		t.Fatalf("unexpected GeneratedAt: %v", r.GeneratedAt)
	}
}

func TestBuild_NoAlerts(t *testing.T) {
	b := digest.NewBuilder(fixedNow)
	r := b.Build([]digest.Entry{
		{Path: "secret/ok", TTL: 72 * time.Hour},
	})
	if r.AlertCount != 0 {
		t.Fatalf("expected 0 alerts, got %d", r.AlertCount)
	}
}

func TestWriteTo(t *testing.T) {
	b := digest.NewBuilder(fixedNow)
	r := b.Build(makeEntries())
	var sb strings.Builder
	r.WriteTo(&sb) //nolint:errcheck
	out := sb.String()
	if !strings.Contains(out, "Total secrets: 3") {
		t.Fatalf("missing total in output: %s", out)
	}
}
