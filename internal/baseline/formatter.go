package baseline

import (
	"fmt"
	"strings"
)

const (
	col1 = 40
	col2 = 14
	col3 = 20
	col4 = 20
)

// FormatDrifts renders a table of Drift values for CLI output.
func FormatDrifts(drifts []Drift) string {
	if len(drifts) == 0 {
		return "No drift detected.\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s %-*s\n", col1, "PATH", col2, "FIELD", col3, "WAS", col4, "NOW")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", col1+col2+col3+col4+3) + "\n")
	for _, d := range drifts {
		sb.WriteString(fmt.Sprintf("%-*s %-*s %-*s %-*s\n",
			col1, truncate(d.Path, col1),
			col2, truncate(d.Field, col2),
			col3, truncate(d.Was, col3),
			col4, truncate(d.Now, col4),
		))
	}
	sb.WriteString(fmt.Sprintf("\n%d drift(s) found.\n", len(drifts)))
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
