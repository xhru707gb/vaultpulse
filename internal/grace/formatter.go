package grace

import (
	"fmt"
	"strings"
	"time"
)

const (
	colPath     = "PATH"
	colExpired  = "EXPIRED AT"
	colGraceEnd = "GRACE ENDS"
	colRemains  = "REMAINING"
)

// FormatTable renders active grace entries as an ASCII table.
func FormatTable(entries []Entry, now time.Time) string {
	if len(entries) == 0 {
		return "No secrets currently in grace period.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s  %-20s  %-20s  %s\n", colPath, colExpired, colGraceEnd, colRemains)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 92))
	for _, e := range entries {
		remaining := e.GraceEndsAt.Sub(now)
		fmt.Fprintf(&sb, "%-40s  %-20s  %-20s  %s\n",
			truncate(e.Path, 40),
			e.ExpiredAt.UTC().Format(time.RFC3339),
			e.GraceEndsAt.UTC().Format(time.RFC3339),
			formatRemaining(remaining),
		)
	}
	return sb.String()
}

func formatRemaining(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
