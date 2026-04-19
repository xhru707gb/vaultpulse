package secretmap

import (
	"fmt"
	"strings"
	"time"
)

const (
	col1 = 40
	col2 = 8
	col3 = 22
	col4 = 16
)

// FormatTable renders entries as a plain-text table.
func FormatTable(entries []Entry, now time.Time) string {
	if len(entries) == 0 {
		return "no secrets registered\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-*s %-*s %-*s %-*s\n",
		col1, "PATH", col2, "VERSION", col3, "EXPIRES", col4, "OWNER")
	sb.WriteString(strings.Repeat("-", col1+col2+col3+col4+3) + "\n")
	for _, e := range entries {
		expiry := formatExpiry(e.ExpiresAt, now)
		fmt.Fprintf(&sb, "%-*s %-*d %-*s %-*s\n",
			col1, truncate(e.Path, col1),
			col2, e.Version,
			col3, expiry,
			col4, truncate(e.Owner, col4),
		)
	}
	return sb.String()
}

func formatExpiry(t time.Time, now time.Time) string {
	if t.IsZero() {
		return "never"
	}
	d := t.Sub(now)
	if d <= 0 {
		return "EXPIRED"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
