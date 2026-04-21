package secretmeta

import (
	"fmt"
	"sort"
	"strings"
)

const (
	maxPathWidth = 36
	maxKeyWidth  = 20
	maxValWidth  = 28
)

// FormatTable renders the metadata registry as an ASCII table.
func FormatTable(r *Registry) string {
	paths := r.Paths()
	if len(paths) == 0 {
		return "No metadata entries found.\n"
	}
	sort.Strings(paths)

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s  %-*s  %s\n", maxPathWidth, "PATH", maxKeyWidth, "KEY", "VALUE")
	fmt.Fprintf(&sb, "%s  %s  %s\n",
		strings.Repeat("-", maxPathWidth),
		strings.Repeat("-", maxKeyWidth),
		strings.Repeat("-", maxValWidth),
	)

	for _, path := range paths {
		meta, err := r.Get(path)
		if err != nil {
			continue
		}
		keys := make([]string, 0, len(meta))
		for k := range meta {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i, k := range keys {
			displayPath := ""
			if i == 0 {
				displayPath = truncate(path, maxPathWidth)
			}
			fmt.Fprintf(&sb, "%-*s  %-*s  %s\n",
				maxPathWidth, displayPath,
				maxKeyWidth, truncate(k, maxKeyWidth),
				truncate(meta[k], maxValWidth),
			)
		}
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
