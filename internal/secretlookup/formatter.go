package secretlookup

import (
	"fmt"
	"sort"
	"strings"
)

const (
	col1Width = 20
	col2Width = 40
)

// FormatTable renders the duplicate fingerprint groups as a plain-text table.
func FormatTable(duplicates map[string][]string) string {
	if len(duplicates) == 0 {
		return "No duplicate secrets detected.\n"
	}

	var sb strings.Builder
	header := fmt.Sprintf("%-*s  %-*s  %s\n",
		col1Width, "FINGERPRINT",
		col2Width, "PATH",
		"SHARED")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", col1Width+col2Width+12) + "\n")

	// Sort fingerprints for deterministic output.
	fps := make([]string, 0, len(duplicates))
	for fp := range duplicates {
		fps = append(fps, fp)
	}
	sort.Strings(fps)

	for _, fp := range fps {
		paths := duplicates[fp]
		sort.Strings(paths)
		shared := len(paths)
		for i, p := range paths {
			if i == 0 {
				fmt.Fprintf(&sb, "%-*s  %-*s  %d\n",
					col1Width, truncate(fp, col1Width),
					col2Width, truncate(p, col2Width),
					shared)
			} else {
				fmt.Fprintf(&sb, "%-*s  %-*s\n",
					col1Width, "",
					col2Width, truncate(p, col2Width))
			}
		}
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}
