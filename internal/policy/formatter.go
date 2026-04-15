package policy

import (
	"fmt"
	"strings"
)

const (
	labelCompliant    = "COMPLIANT"
	labelViolation    = "VIOLATION"
	colWidthPath      = 30
	colWidthPolicy    = 16
	colWidthStatus    = 12
)

// FormatTable renders policy check results as an ASCII table.
func FormatTable(statuses []Status) string {
	if len(statuses) == 0 {
		return "No policy results to display.\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s %s\n",
		colWidthPath, "PATH",
		colWidthPolicy, "POLICY",
		colWidthStatus, "STATUS",
		"VIOLATIONS")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")
	for _, s := range statuses {
		label := complianceLabel(s.Compliant)
		violations := strings.Join(s.Violations, "; ")
		if violations == "" {
			violations = "-"
		}
		sb.WriteString(fmt.Sprintf("%-*s %-*s %-*s %s\n",
			colWidthPath, truncatePath(s.Path, colWidthPath),
			colWidthPolicy, s.Policy,
			colWidthStatus, label,
			violations))
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of policy evaluation results.
func FormatSummary(statuses []Status) string {
	total := len(statuses)
	violations := 0
	for _, s := range statuses {
		if !s.Compliant {
			violations++
		}
	}
	return fmt.Sprintf("Policy check: %d evaluated, %d compliant, %d violations\n",
		total, total-violations, violations)
}

func complianceLabel(compliant bool) string {
	if compliant {
		return labelCompliant
	}
	return labelViolation
}

func truncatePath(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
