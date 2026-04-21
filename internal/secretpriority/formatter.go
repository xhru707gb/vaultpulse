package secretpriority

import (
	"fmt"
	"strings"
)

const (
	maxPathWidth = 40
	maxRuleWidth = 20
)

// FormatTable renders a slice of Results as an ASCII table.
func FormatTable(results []Result) string {
	if len(results) == 0 {
		return "No priority results to display.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-42s %-10s %s\n", "PATH", "PRIORITY", "MATCHED RULE")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 72))

	for _, r := range results {
		fmt.Fprintf(&sb, "%-42s %-10s %s\n",
			truncate(r.Path, maxPathWidth),
			LevelLabel(r.Level),
			truncate(r.Rule, maxRuleWidth),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line count summary grouped by level.
func FormatSummary(results []Result) string {
	counts := map[Level]int{}
	for _, r := range results {
		counts[r.Level]++
	}
	return fmt.Sprintf(
		"Total: %d | Critical: %d | High: %d | Medium: %d | Low: %d\n",
		len(results),
		counts[LevelCritical],
		counts[LevelHigh],
		counts[LevelMedium],
		counts[LevelLow],
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
