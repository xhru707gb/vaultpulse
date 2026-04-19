package secretgroup

import (
	"fmt"
	"strings"
)

const (
	maxPathLen = 48
	colWidth   = 20
)

// FormatTable renders all groups as a plain-text table.
func FormatTable(groups []*Group) string {
	if len(groups) == 0 {
		return "no groups defined\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s  %-6s  %s\n", colWidth, "GROUP", "PATHS", "SAMPLE PATH")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", colWidth+2+6+2+maxPathLen))
	for _, g := range groups {
		sample := ""
		if len(g.Paths) > 0 {
			sample = truncate(g.Paths[0], maxPathLen)
		}
		fmt.Fprintf(&sb, "%-*s  %-6d  %s\n", colWidth, truncate(g.Name, colWidth), len(g.Paths), sample)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of group counts.
func FormatSummary(groups []*Group) string {
	total := 0
	for _, g := range groups {
		total += len(g.Paths)
	}
	return fmt.Sprintf("%d group(s), %d total path(s)\n", len(groups), total)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
