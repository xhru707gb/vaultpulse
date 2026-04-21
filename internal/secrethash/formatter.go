package secrethash

import (
	"fmt"
	"strings"
	"time"
)

const (
	maxPathLen = 40
	maxHashLen = 16
)

// FormatTable renders tracked hash entries as a plain-text table.
func FormatTable(entries []Entry) string {
	if len(entries) == 0 {
		return "No hash entries tracked.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-42s %-18s %-8s %s\n",
		"PATH", "HASH (PREFIX)", "VERSION", "CHANGED AT")
	fmt.Fprintln(&sb, strings.Repeat("-", 85))

	for _, e := range entries {
		fmt.Fprintf(&sb, "%-42s %-18s %-8d %s\n",
			truncate(e.Path, maxPathLen),
			truncateHash(e.Hash, maxHashLen),
			e.Version,
			formatTime(e.ChangedAt),
		)
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of tracked hashes.
func FormatSummary(entries []Entry) string {
	return fmt.Sprintf("Tracked secrets: %d\n", len(entries))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func truncateHash(h string, max int) string {
	if len(h) <= max {
		return h
	}
	return h[:max] + "…"
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}
