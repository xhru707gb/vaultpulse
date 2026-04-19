package secretexpiry

import (
	"fmt"
	"strings"
	"time"
)

const (
	labelExpired = "EXPIRED"
	labelWarning = "WARNING"
	labelOK      = "OK"
)

func stateLabel(s Status) string {
	switch {
	case s.Expired:
		return labelExpired
	case s.Warning:
		return labelWarning
	default:
		return labelOK
	}
}

func formatRemaining(d time.Duration) string {
	if d <= 0 {
		return "expired"
	}
	h := int(d.Hours())
	if h >= 24 {
		return fmt.Sprintf("%dd", h/24)
	}
	return fmt.Sprintf("%dh%dm", h, int(d.Minutes())%60)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "..." + s[len(s)-n+3:]
}

// FormatTable renders statuses as a plain-text table.
func FormatTable(statuses []Status) string {
	if len(statuses) == 0 {
		return "no secret expiry entries\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s %-10s %s\n", "PATH", "STATE", "REMAINING")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 64))
	for _, s := range statuses {
		fmt.Fprintf(&sb, "%-40s %-10s %s\n",
			truncate(s.Path, 40),
			stateLabel(s),
			formatRemaining(s.Remaining),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of expiry statuses.
func FormatSummary(statuses []Status) string {
	var expired, warning, ok int
	for _, s := range statuses {
		switch {
		case s.Expired:
			expired++
		case s.Warning:
			warning++
		default:
			ok++
		}
	}
	return fmt.Sprintf("total=%d expired=%d warning=%d ok=%d",
		len(statuses), expired, warning, ok)
}
