package secretclassify

import (
	"fmt"
	"strings"
)

const (
	colWidth = 40
	lvlWidth = 14
)

// FormatTable renders classification results as a plain-text table.
func FormatTable(results []Result) string {
	if len(results) == 0 {
		return "no classification results\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s\n", colWidth, "PATH", lvlWidth, "LEVEL")
	sep := strings.Repeat("-", colWidth+lvlWidth+1)
	sb.WriteString(header)
	sb.WriteString(sep + "\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-*s %-*s\n",
			colWidth, truncate(r.Path, colWidth),
			lvlWidth, levelLabel(r.Level),
		))
	}
	return sb.String()
}

// FormatSummary returns a one-line count summary grouped by level.
func FormatSummary(results []Result) string {
	counts := map[Level]int{}
	for _, r := range results {
		counts[r.Level]++
	}
	return fmt.Sprintf("total=%d  secret=%d  confidential=%d  internal=%d  public=%d\n",
		len(results),
		counts[LevelSecret],
		counts[LevelConfidential],
		counts[LevelInternal],
		counts[LevelPublic],
	)
}

func levelLabel(l Level) string {
	switch l {
	case LevelSecret:
		return "SECRET"
	case LevelConfidential:
		return "CONFIDENTIAL"
	case LevelInternal:
		return "INTERNAL"
	default:
		return "PUBLIC"
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
