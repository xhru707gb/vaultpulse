package escalation

import (
	"fmt"
	"strings"
	"time"
)

// FormatTable renders escalation events as a plain-text table.
func FormatTable(events []Event) string {
	if len(events) == 0 {
		return "No escalation events.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s %-10s %s\n", "PATH", "LEVEL", "TTL")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 62))
	for _, ev := range events {
		fmt.Fprintf(&sb, "%-40s %-10s %s\n",
			truncatePath(ev.Path, 40),
			levelLabel(ev.Level),
			formatTTL(ev.TTL),
		)
	}
	return sb.String()
}

func levelLabel(l Level) string {
	switch l {
	case LevelCritical:
		return "CRITICAL"
	case LevelWarning:
		return "WARNING"
	default:
		return "INFO"
	}
}

func formatTTL(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

func truncatePath(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return "..." + s[len(s)-(max-3):]
}
