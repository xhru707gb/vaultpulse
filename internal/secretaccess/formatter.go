package secretaccess

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxPathWidth = 40
	maxTimeWidth = 20
)

// FormatTable renders access tracking entries as a plain-text table.
func FormatTable(entries []AccessEntry) string {
	if len(entries) == 0 {
		return "no access records found\n"
	}

	var sb strings.Builder

	header := fmt.Sprintf("%-42s  %-6s  %-22s  %s\n",
		"PATH", "COUNT", "LAST ACCESS", "FIRST ACCESS")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", 90) + "\n")

	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-42s  %-6d  %-22s  %s\n",
			truncate(e.Path, maxPathWidth),
			e.Count,
			formatTime(e.LastAccess),
			formatTime(e.FirstAccess),
		))
	}

	return sb.String()
}

// FormatSummary returns a one-line summary of access statistics.
func FormatSummary(entries []AccessEntry) string {
	if len(entries) == 0 {
		return "total paths: 0, total accesses: 0\n"
	}
	var total int
	for _, e := range entries {
		total += e.Count
	}
	return fmt.Sprintf("total paths: %d, total accesses: %d\n", len(entries), total)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
