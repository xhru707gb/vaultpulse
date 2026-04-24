package secretreview

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxPathLen = 40
	headerSep  = "-"
)

// FormatTable renders review entries as a plain-text table.
func FormatTable(entries []*Entry) string {
	if len(entries) == 0 {
		return "no review entries found\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-42s %-16s %-10s %-22s %s\n",
		"PATH", "REVIEWER", "STATUS", "NEXT REVIEW", "INTERVAL")
	sb.WriteString(strings.Repeat(headerSep, 100) + "\n")
	for _, e := range entries {
		fmt.Fprintf(&sb, "%-42s %-16s %-10s %-22s %s\n",
			truncate(e.Path, maxPathLen),
			truncate(e.Reviewer, 14),
			statusLabel(e.Status),
			e.NextReview.UTC().Format(time.RFC3339),
			formatInterval(e.Interval),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of review statuses.
func FormatSummary(entries []*Entry) string {
	var overdue, ok int
	for _, e := range entries {
		if e.Status == StatusOverdue {
			overdue++
		} else {
			ok++
		}
	}
	return fmt.Sprintf("total=%d approved=%d overdue=%d", len(entries), ok, overdue)
}

func statusLabel(s Status) string {
	switch s {
	case StatusApproved:
		return "approved"
	case StatusOverdue:
		return "overdue"
	default:
		return "pending"
	}
}

func formatInterval(d time.Duration) string {
	h := int(d.Hours())
	if h >= 24 {
		return fmt.Sprintf("%dd", h/24)
	}
	return fmt.Sprintf("%dh", h)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
