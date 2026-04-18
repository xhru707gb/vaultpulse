package tagindex

import (
	"fmt"
	"strings"
)

const (
	colTag   = "TAG"
	colCount = "PATHS"
	colList  = "MEMBERS"
)

// FormatTable renders the index as a plain-text table.
func FormatTable(idx *Index) string {
	tags := idx.Tags()
	if len(tags) == 0 {
		return "(no tags registered)\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-24s %6s  %s\n", colTag, colCount, colList)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 72))

	for _, tag := range tags {
		paths := idx.Paths(tag)
		preview := truncatePaths(paths, 3)
		fmt.Fprintf(&sb, "%-24s %6d  %s\n", truncate(tag, 24), len(paths), preview)
	}
	return sb.String()
}

func truncatePaths(paths []string, max int) string {
	if len(paths) <= max {
		return strings.Join(paths, ", ")
	}
	preview := paths[:max]
	return strings.Join(preview, ", ") + fmt.Sprintf(" (+%d more)", len(paths)-max)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
