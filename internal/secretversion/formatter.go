package secretversion

import (
	"fmt"
	"strings"
	"time"
)

const maxPathLen = 48

// FormatTable renders version entries as a plain-text table.
func FormatTable(entries []Entry) string {
	if len(entries) == 0 {
		return "no secret versions tracked\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-50s  %7s  %-20s  %-20s\n",
		"PATH", "VERSION", "CREATED", "UPDATED")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 102))
	for _, e := range entries {
		fmt.Fprintf(&sb, "%-50s  %7d  %-20s  %-20s\n",
			truncate(e.Path, maxPathLen),
			e.Version,
			formatTime(e.CreatedAt),
			formatTime(e.UpdatedAt),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary.
func FormatSummary(entries []Entry) string {
	if len(entries) == 0 {
		return "tracked secrets: 0\n"
	}
	max := 0
	for _, e := range entries {
		if e.Version > max {
			max = e.Version
		}
	}
	return fmt.Sprintf("tracked secrets: %d  highest version: %d\n", len(entries), max)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.UTC().Format("2006-01-02 15:04")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
