package quarantine

import (
	"fmt"
	"strings"
	"time"
)

const (
	col1 = 40
	col2 = 18
	col3 = 24
	col4 = 20
)

// FormatTable renders quarantined entries as a plain-text table.
func FormatTable(entries []Entry) string {
	var b strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s %s",
		col1, "PATH",
		col2, "REASON",
		col3, "QUARANTINED AT",
		"NOTE",
	)
	b.WriteString(header + "\n")
	b.WriteString(strings.Repeat("-", len(header)+10) + "\n")

	if len(entries) == 0 {
		b.WriteString("  no quarantined secrets\n")
		return b.String()
	}

	for _, e := range entries {
		line := fmt.Sprintf("%-*s %-*s %-*s %s",
			col1, truncate(e.Path, col1),
			col2, string(e.Reason),
			col3, e.QuarantinedAt.UTC().Format(time.RFC3339),
			e.Note,
		)
		b.WriteString(line + "\n")
	}
	return b.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
