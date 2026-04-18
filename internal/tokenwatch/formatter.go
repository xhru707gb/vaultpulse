package tokenwatch

import (
	"fmt"
	"strings"
	"time"
)

// FormatTable renders token statuses as a plain-text table.
func FormatTable(statuses []Status) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-20s %-24s %-10s %s\n", "ACCESSOR", "DISPLAY NAME", "STATE", "REMAINING")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 72))
	for _, s := range statuses {
		fmt.Fprintf(&sb, "%-20s %-24s %-10s %s\n",
			truncate(s.Token.Accessor, 20),
			truncate(s.Token.DisplayName, 24),
			stateLabel(s.State),
			formatRemaining(s.Remaining),
		)
	}
	return sb.String()
}

func stateLabel(state string) string {
	switch state {
	case "expired":
		return "EXPIRED"
	case "warning":
		return "WARNING"
	default:
		return "OK"
	}
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
