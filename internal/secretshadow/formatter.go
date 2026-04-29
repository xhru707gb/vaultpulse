package secretshadow

import (
	"fmt"
	"strings"
	"time"
)

const (
	labelOK      = "OK"
	labelDiverged = "DIVERGED"
	maxPathLen   = 40
	maxHashLen   = 16
)

func divergedLabel(d bool) string {
	if d {
		return labelDiverged
	}
	return labelOK
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func formatCaptured(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.UTC().Format(time.RFC3339)
}

// FormatTable renders shadow entries as a plain-text table.
func FormatTable(entries []Entry) string {
	if len(entries) == 0 {
		return "no shadow entries recorded\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-42s %-18s %-18s %s\n",
		"PATH", "HASH (prefix)", "CAPTURED AT", "STATUS")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 90))
	for _, e := range entries {
		hashPrefix := truncate(e.Hash, maxHashLen)
		fmt.Fprintf(&sb, "%-42s %-18s %-18s %s\n",
			truncate(e.Path, maxPathLen),
			hashPrefix,
			formatCaptured(e.CapturedAt),
			divergedLabel(e.Diverged),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of shadow status.
func FormatSummary(entries []Entry) string {
	total := len(entries)
	diverged := 0
	for _, e := range entries {
		if e.Diverged {
			diverged++
		}
	}
	return fmt.Sprintf("shadow: %d total, %d diverged, %d ok\n",
		total, diverged, total-diverged)
}
