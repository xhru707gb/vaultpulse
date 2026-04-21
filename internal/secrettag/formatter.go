package secrettag

import (
	"fmt"
	"strings"
)

const (
	maxPathWidth = 40
	maxTagsWidth = 50
)

// FormatTable renders a human-readable table of path→tags mappings.
// paths is a sorted slice of paths; for each path Tags() is called.
func FormatTable(t *Tagger, paths []string) string {
	if len(paths) == 0 {
		return "No tagged secrets.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s  %s\n", maxPathWidth, "PATH", "TAGS")
	fmt.Fprintf(&sb, "%s  %s\n", strings.Repeat("-", maxPathWidth), strings.Repeat("-", maxTagsWidth))

	for _, path := range paths {
		tags, err := t.Tags(path)
		if err != nil {
			continue
		}
		tagLine := strings.Join(tags, ", ")
		fmt.Fprintf(&sb, "%-*s  %s\n",
			maxPathWidth, truncate(path, maxPathWidth),
			truncate(tagLine, maxTagsWidth),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of tagging statistics.
func FormatSummary(t *Tagger) string {
	return fmt.Sprintf("Tagged paths: %d\n", t.Len())
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
