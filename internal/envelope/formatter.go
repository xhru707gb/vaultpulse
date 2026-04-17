package envelope

import (
	"fmt"
	"strings"
	"time"
)

const (
	colPath    = "PATH"
	colKey     = "KEY VERSION"
	colAge     = "AGE"
	colEncAt   = "ENCRYPTED AT"
)

// FormatTable renders a slice of Envelopes as a plain-text table.
func FormatTable(envelopes []*Envelope) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%-40s %-20s %-12s %s\n", colPath, colKey, colAge, colEncAt)
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 90))
	for _, e := range envelopes {
		age := formatAge(e.Age())
		encAt := e.EncryptedAt.Format(time.RFC3339)
		path := truncate(e.Path, 38)
		fmt.Fprintf(&sb, "%-40s %-20s %-12s %s\n", path, e.KeyVersion, age, encAt)
	}
	return sb.String()
}

func formatAge(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
