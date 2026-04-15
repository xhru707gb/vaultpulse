package diff

import (
	"fmt"
	"strings"
	"time"
)

// FormatTable renders a slice of ChangeEntry values as a human-readable table.
func FormatTable(changes []ChangeEntry) string {
	if len(changes) == 0 {
		return "No changes detected.\n"
	}

	var sb strings.Builder

	header := fmt.Sprintf("%-40s %-12s %-26s %-26s\n",
		"PATH", "CHANGE", "OLD VALUE", "NEW VALUE")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", 108) + "\n")

	for _, c := range changes {
		old := formatField(c.OldValue)
		new := formatField(c.NewValue)
		line := fmt.Sprintf("%-40s %-12s %-26s %-26s\n",
			truncatePath(c.Path, 40),
			changeLabel(c.Kind),
			old,
			new,
		)
		sb.WriteString(line)
	}

	return sb.String()
}

// FormatSummary returns a one-line summary of the diff result.
func FormatSummary(changes []ChangeEntry) string {
	added, removed, modified := 0, 0, 0
	for _, c := range changes {
		switch c.Kind {
		case KindAdded:
			added++
		case KindRemoved:
			removed++
		case KindModified:
			modified++
		}
	}
	return fmt.Sprintf("Diff summary: +%d added, -%d removed, ~%d modified", added, removed, modified)
}

func changeLabel(k ChangeKind) string {
	switch k {
	case KindAdded:
		return "ADDED"
	case KindRemoved:
		return "REMOVED"
	case KindModified:
		return "MODIFIED"
	default:
		return "UNKNOWN"
	}
}

func formatField(v interface{}) string {
	if v == nil {
		return "-"
	}
	switch val := v.(type) {
	case string:
		if val == "" {
			return "-"
		}
		return val
	case time.Time:
		if val.IsZero() {
			return "-"
		}
		return val.UTC().Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func truncatePath(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
