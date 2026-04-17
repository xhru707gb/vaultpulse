package fingerprint

import (
	"fmt"
	"strings"
)

const (
	col1 = 40
	col2 = 16
	col3 = 8
)

// FormatTable renders a slice of Results as a plain-text table.
func FormatTable(results []Result) string {
	if len(results) == 0 {
		return "No fingerprint results.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s  %-*s  %s\n", col1, "PATH", col2, "FINGERPRINT", "KEYS")
	fmt.Fprintf(&sb, "%s  %s  %s\n",
		strings.Repeat("-", col1),
		strings.Repeat("-", col2),
		strings.Repeat("-", col3))

	for _, r := range results {
		short := r.Fingerprint
		if len(short) > col2 {
			short = short[:col2]
		}
		fmt.Fprintf(&sb, "%-*s  %-*s  %d\n",
			col1, truncate(r.Path, col1),
			col2, short,
			r.KeyCount)
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
