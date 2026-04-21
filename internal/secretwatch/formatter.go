package secretwatch

import (
	"fmt"
	"strings"
)

const (
	_pathWidth  = 40
	_kindWidth  = 10
	_detailWidth = 30
)

// FormatTable renders a slice of Events as a human-readable table.
func FormatTable(events []Event) string {
	if len(events) == 0 {
		return "No changes detected.\n"
	}
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s\n",
		_pathWidth, "PATH",
		_kindWidth, "CHANGE",
		_detailWidth, "DETAIL")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", _pathWidth+_kindWidth+_detailWidth+2) + "\n")
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%-*s %-*s %-*s\n",
			_pathWidth, truncate(e.Path, _pathWidth),
			_kindWidth, kindLabel(e.Kind),
			_detailWidth, truncate(e.Detail, _detailWidth),
		))
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of the events slice.
func FormatSummary(events []Event) string {
	var added, removed, modified int
	for _, e := range events {
		switch e.Kind {
		case "added":
			added++
		case "removed":
			removed++
		case "modified":
			modified++
		}
	}
	return fmt.Sprintf("added=%d removed=%d modified=%d", added, removed, modified)
}

func kindLabel(kind string) string {
	switch kind {
	case "added":
		return "ADDED"
	case "removed":
		return "REMOVED"
	case "modified":
		return "MODIFIED"
	default:
		return strings.ToUpper(kind)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
