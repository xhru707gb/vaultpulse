package secretdrift

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxPathWidth = 40
	maxHashWidth = 16
)

// FormatTable renders drift entries as an ASCII table.
func FormatTable(drifts []DriftEntry) string {
	if len(drifts) == 0 {
		return "No drift detected.\n"
	}

	var sb strings.Builder
	header := fmt.Sprintf("%-*s  %-*s  %-*s  %s\n",
		maxPathWidth, "PATH",
		maxHashWidth, "PREV HASH",
		maxHashWidth, "CURR HASH",
		"DETECTED AT")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", len(header)-1) + "\n")

	for _, d := range drifts {
		sb.WriteString(fmt.Sprintf("%-*s  %-*s  %-*s  %s\n",
			maxPathWidth, truncate(d.Path, maxPathWidth),
			maxHashWidth, truncate(d.PreviousHash, maxHashWidth),
			maxHashWidth, truncate(d.CurrentHash, maxHashWidth),
			d.DetectedAt.UTC().Format(time.RFC3339),
		))
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of drift counts.
func FormatSummary(drifts []DriftEntry) string {
	return fmt.Sprintf("Drift summary: %d change(s) detected.\n", len(drifts))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
