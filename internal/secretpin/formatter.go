package secretpin

import (
	"fmt"
	"strings"
)

const (
	colWidthPath   = 36
	colWidthVer    = 8
	colWidthBy     = 16
	colWidthDrift  = 8
)

// FormatPins renders a table of active pins.
func FormatPins(pins []Pin) string {
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s",
		colWidthPath, "PATH",
		colWidthVer, "VERSION",
		colWidthBy, "PINNED BY",
	)
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")
	for _, p := range pins {
		sb.WriteString(fmt.Sprintf("%-*s %-*d %-*s\n",
			colWidthPath, truncate(p.Path, colWidthPath),
			colWidthVer, p.Version,
			colWidthBy, truncate(p.PinnedBy, colWidthBy),
		))
	}
	return sb.String()
}

// FormatDrifts renders a table of drift check results.
func FormatDrifts(results []DriftResult) string {
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s %-*s",
		colWidthPath, "PATH",
		colWidthVer, "PINNED",
		colWidthVer, "CURRENT",
		colWidthDrift, "DRIFTED",
	)
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", len(header)) + "\n")
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("%-*s %-*d %-*d %-*s\n",
			colWidthPath, truncate(r.Path, colWidthPath),
			colWidthVer, r.PinnedVersion,
			colWidthVer, r.CurrentVersion,
			colWidthDrift, driftLabel(r.Drifted),
		))
	}
	return sb.String()
}

func driftLabel(drifted bool) string {
	if drifted {
		return "YES"
	}
	return "no"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
