package digest

import (
	"fmt"
	"strings"
	"time"
)

const (
	colPath  = 40
	colState = 14
	colTTL   = 12
)

// FormatTable renders a Report as a fixed-width table string.
func FormatTable(r Report) string {
	var sb strings.Builder
	header := fmt.Sprintf("%-*s %-*s %-*s\n",
		colPath, "PATH", colState, "STATE", colTTL, "TTL")
	sb.WriteString(header)
	sb.WriteString(strings.Repeat("-", colPath+colState+colTTL+2) + "\n")
	for _, e := range r.Entries {
		path := truncatePath(e.Path, colPath)
		state := stateLabel(e)
		ttl := formatTTL(e.TTL)
		sb.WriteString(fmt.Sprintf("%-*s %-*s %-*s\n",
			colPath, path, colState, state, colTTL, ttl))
	}
	sb.WriteString(fmt.Sprintf("\nGenerated: %s  Total: %d  Alerts: %d\n",
		r.GeneratedAt.Format(time.RFC3339), r.TotalSecrets, r.AlertCount))
	return sb.String()
}

func stateLabel(e Entry) string {
	switch {
	case e.Expired:
		return "EXPIRED"
	case e.ExpiresSoon:
		return "EXPIRING SOON"
	case e.Overdue:
		return "OVERDUE"
	case e.Unhealthy:
		return "UNHEALTHY"
	default:
		return "OK"
	}
}

func formatTTL(d time.Duration) string {
	if d <= 0 {
		return "n/a"
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%02dm", h, m)
}

func truncatePath(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return "..." + s[len(s)-(max-3):]
}
