package secretage

import (
	"fmt"
	"strings"
	"time"
)

const (
	labelOverdue = "OVERDUE"
	labelOK      = "OK"
)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func formatAge(d time.Duration) string {
	days := int(d.Hours()) / 24
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	h := int(d.Hours())
	if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}

// FormatTable renders secret age statuses as a plain-text table.
func FormatTable(statuses []Status) string {
	if len(statuses) == 0 {
		return "no secrets tracked\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s  %-8s  %-8s  %s\n", "PATH", "AGE", "MAX AGE", "STATUS")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 72))
	for _, s := range statuses {
		label := labelOK
		if s.Overdue {
			label = labelOverdue
		}
		fmt.Fprintf(&sb, "%-40s  %-8s  %-8s  %s\n",
			truncate(s.Path, 40),
			formatAge(s.Age),
			formatAge(s.MaxAge),
			label,
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary.
func FormatSummary(statuses []Status) string {
	overdue := 0
	for _, s := range statuses {
		if s.Overdue {
			overdue++
		}
	}
	return fmt.Sprintf("total=%d overdue=%d\n", len(statuses), overdue)
}
