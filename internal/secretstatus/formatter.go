package secretstatus

import (
	"fmt"
	"strings"
)

const (
	colWidth = 60
)

// FormatTable renders a slice of Entry values as a plain-text table.
func FormatTable(entries []*Entry) string {
	if len(entries) == 0 {
		return "No secret status entries.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s  %-10s  %s\n", "PATH", "STATUS", "REASONS")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", colWidth+20))
	for _, e := range entries {
		reasons := strings.Join(e.Reasons, "; ")
		if reasons == "" {
			reasons = "-"
		}
		fmt.Fprintf(&sb, "%-40s  %-10s  %s\n",
			truncate(e.Path, 40),
			levelLabel(e.Level),
			reasons,
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of entry counts by level.
func FormatSummary(entries []*Entry) string {
	ok, warn, crit := 0, 0, 0
	for _, e := range entries {
		switch e.Level {
		case LevelOK:
			ok++
		case LevelWarning:
			warn++
		case LevelCritical:
			crit++
		}
	}
	return fmt.Sprintf("Total: %d  OK: %d  Warning: %d  Critical: %d",
		len(entries), ok, warn, crit)
}

func levelLabel(l Level) string {
	switch l {
	case LevelOK:
		return "OK"
	case LevelWarning:
		return "WARNING"
	case LevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
