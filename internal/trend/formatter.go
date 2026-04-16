package trend

import (
	"fmt"
	"strings"
)

const (
	barMax  = 20
	headerFmt = "%-40s %-12s %6s  %s\n"
	rowFmt   = "%-40s %-12s %6d  %s\n"
)

// FormatTable renders trend reports as an ASCII table with spark-bar.
func FormatTable(reports []Report) string {
	if len(reports) == 0 {
		return "No trend data available.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(headerFmt, "PATH", "EVENT", "TOTAL", "TREND"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, r := range reports {
		bar := sparkBar(r.Points)
		sb.WriteString(fmt.Sprintf(rowFmt, truncate(r.Path, 40), r.EventType, r.TotalCount, bar))
	}
	return sb.String()
}

func sparkBar(points []Point) string {
	if len(points) == 0 {
		return ""
	}
	max := 0
	for _, p := range points {
		if p.Count > max {
			max = p.Count
		}
	}
	blocks := []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var sb strings.Builder
	for _, p := range points {
		idx := 0
		if max > 0 {
			idx = int(float64(p.Count) / float64(max) * float64(len(blocks)-1))
		}
		sb.WriteString(blocks[idx])
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
