package ttlpolicy

import (
	"fmt"
	"strings"
)

const (
	colWidth = 40
	ttlWidth = 12
)

// FormatTable renders TTL policy results as a plain-text table.
func FormatTable(results []Result) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s %-*s %-10s %s\n", colWidth, "PATH", ttlWidth, "TTL", "STATUS", "VIOLATION")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", colWidth+ttlWidth+30))
	for _, r := range results {
		label := complianceLabel(r.Compliant)
		fmt.Fprintf(&sb, "%-*s %-*s %-10s %s\n",
			colWidth, truncate(r.Path, colWidth),
			ttlWidth, r.TTL.String(),
			label,
			r.Violation,
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of compliance results.
func FormatSummary(results []Result) string {
	total := len(results)
	violations := 0
	for _, r := range results {
		if !r.Compliant {
			violations++
		}
	}
	if violations == 0 {
		return fmt.Sprintf("All %d secret(s) comply with TTL policy.", total)
	}
	return fmt.Sprintf("%d/%d secret(s) violate TTL policy.", violations, total)
}

func complianceLabel(ok bool) string {
	if ok {
		return "COMPLIANT"
	}
	return "VIOLATION"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
