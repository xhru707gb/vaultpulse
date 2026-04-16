package metrics

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	colPath      = "PATH"
	colStatus    = "STATUS"
	colCheckedAt = "CHECKED AT"
	colDuration  = "DURATION"
	colError     = "ERROR"
)

// FormatTable renders snapshots as a fixed-width ASCII table.
func FormatTable(snapshots []Snapshot) string {
	if len(snapshots) == 0 {
		return "no metrics recorded\n"
	}

	sorted := make([]Snapshot, len(snapshots))
	copy(sorted, snapshots)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Path < sorted[j].Path
	})

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s %-10s %-22s %-12s %s\n",
		colPath, colStatus, colCheckedAt, colDuration, colError)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 100))

	for _, s := range sorted {
		fmt.Fprintf(&sb, "%-40s %-10s %-22s %-12s %s\n",
			truncate(s.Path, 40),
			s.Status,
			formatTime(s.CheckedAt),
			formatDuration(s.Duration),
			s.Error,
		)
	}
	return sb.String()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	return d.Round(time.Millisecond).String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
