package rollup

import (
	"fmt"
	"strings"
)

// FormatTable renders a Summary as a plain-text table.
func FormatTable(s Summary) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Rollup Summary  window=%-10s  total=%d  flushed=%s\n",
		s.Window, s.Total, s.FlushedAt.UTC().Format("2006-01-02T15:04:05Z")))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	sb.WriteString(fmt.Sprintf("%-40s  %-10s  %s\n", "PATH", "LEVEL", "MESSAGE"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")

	for _, e := range s.Events {
		path := truncate(e.Path, 38)
		sb.WriteString(fmt.Sprintf("%-40s  %-10s  %s\n", path, e.Level, e.Message))
	}

	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, lvl := range []string{"ok", "warning", "expired"} {
		if c, ok := s.ByLevel[lvl]; ok {
			sb.WriteString(fmt.Sprintf("  %-10s %d\n", lvl, c))
		}
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
