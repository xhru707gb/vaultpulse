package secretscore

import (
	"fmt"
	"strings"
)

const (
	_pathWidth  = 40
	_reasonWidth = 35
)

// FormatTable renders a slice of Results as a plain-text table.
func FormatTable(results []Result) string {
	if len(results) == 0 {
		return "no results\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s  %6s  %-10s  %s\n", _pathWidth, "PATH", "SCORE", "RISK", "REASON")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", _pathWidth+2+6+2+10+2+_reasonWidth))
	for _, r := range results {
		fmt.Fprintf(&sb, "%-*s  %6d  %-10s  %s\n",
			_pathWidth, truncate(r.Path, _pathWidth),
			r.Score,
			r.Level,
			truncate(r.Reason, _reasonWidth),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of results.
func FormatSummary(results []Result) string {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Level]++
	}
	return fmt.Sprintf("total=%d critical=%d high=%d medium=%d low=%d",
		len(results),
		counts[RiskCritical],
		counts[RiskHigh],
		counts[RiskMedium],
		counts[RiskLow],
	)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
