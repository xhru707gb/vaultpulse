package redact

import (
	"fmt"
	"strings"
)

// FormatTable renders a redacted key/value map as a simple two-column
// plain-text table suitable for CLI output.
func FormatTable(m map[string]string) string {
	if len(m) == 0 {
		return "(no fields)\n"
	}

	const colWidth = 24
	var sb strings.Builder

	header := fmt.Sprintf("%-*s  %s\n", colWidth, "FIELD", "VALUE")
	sep := strings.Repeat("-", colWidth) + "  " + strings.Repeat("-", 40) + "\n"
	sb.WriteString(header)
	sb.WriteString(sep)

	for k, v := range m {
		sb.WriteString(fmt.Sprintf("%-*s  %s\n", colWidth, truncate(k, colWidth), v))
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
