package secretnamespace

import (
	"fmt"
	"strings"
)

const (
	colNS    = 24
	colCount = 8
	colPaths = 40
)

// FormatTable renders the namespace registry as a plain-text table.
func FormatTable(r *Registry) string {
	var sb strings.Builder

	header := fmt.Sprintf("%-*s %-*s %-*s",
		colNS, "NAMESPACE",
		colCount, "PATHS",
		colPaths, "SAMPLE PATH",
	)
	sep := strings.Repeat("-", len(header))
	sb.WriteString(header + "\n")
	sb.WriteString(sep + "\n")

	for _, ns := range r.Namespaces() {
		paths, _ := r.Paths(ns)
		sample := ""
		if len(paths) > 0 {
			sample = truncate(paths[0], colPaths)
		}
		line := fmt.Sprintf("%-*s %-*d %-*s",
			colNS, truncate(ns, colNS),
			colCount, len(paths),
			colPaths, sample,
		)
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of the registry.
func FormatSummary(r *Registry) string {
	ns := r.Namespaces()
	total := 0
	for _, n := range ns {
		paths, _ := r.Paths(n)
		total += len(paths)
	}
	return fmt.Sprintf("namespaces=%d total_paths=%d", len(ns), total)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
