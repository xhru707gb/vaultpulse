package window

import (
	"fmt"
	"strings"
	"time"
)

const (
	colWidth = 30
	timeLayout = "2006-01-02 15:04:05"
)

// FormatTable renders window entries as a plain-text table.
// fmtValue converts a generic value to its display string.
func FormatTable[T any](entries []Entry[T], fmtValue func(T) string) string {
	var sb strings.Builder
	header := fmt.Sprintf("%-4s  %-*s  %s\n", "#", colWidth, "VALUE", "RECORDED AT")
	sep := strings.Repeat("-", len(header)-1) + "\n"
	sb.WriteString(header)
	sb.WriteString(sep)
	if len(entries) == 0 {
		sb.WriteString("no entries in window\n")
		return sb.String()
	}
	for i, e := range entries {
		v := fmtValue(e.Value)
		if len(v) > colWidth {
			v = v[:colWidth-1] + "…"
		}
		sb.WriteString(fmt.Sprintf("%-4d  %-*s  %s\n", i+1, colWidth, v, e.RecordedAt.UTC().Format(timeLayout)))
	}
	return sb.String()
}

// FormatSummary returns a one-line summary of the window contents.
func FormatSummary[T any](entries []Entry[T], duration time.Duration) string {
	return fmt.Sprintf("window(%s): %d entries", duration.String(), len(entries))
}
